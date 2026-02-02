package domain

import (
	"context"
	"database/sql"
	"strings"
	"time"

	pkgerr "github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/queries"

	"github.com/Bnei-Baruch/gxydb-api/common"
)

type RoomServerAssignment struct {
	RoomID      int64
	GatewayName string
	Region      string
	AssignedAt  time.Time
	LastUsedAt  time.Time
}

type ServerLoad struct {
	GatewayName string
	Load        int
}

// GatewayStatusChecker interface for checking if gateway is online
type GatewayStatusChecker interface {
	IsGatewayOnline(serverName string) bool
}

type RoomServerAssignmentManager struct {
	db                 common.DBInterface
	availableServers   []string
	maxServerCapacity  int
	avgRoomOccupancy   int
	serverRegions      map[string][]string
	statusChecker      GatewayStatusChecker // Optional: for checking online status
}

func NewRoomServerAssignmentManager(db common.DBInterface, servers []string, maxCapacity, avgOccupancy int, regions map[string][]string) *RoomServerAssignmentManager {
	return &RoomServerAssignmentManager{
		db:                 db,
		availableServers:   servers,
		maxServerCapacity:  maxCapacity,
		avgRoomOccupancy:   avgOccupancy,
		serverRegions:      regions,
		statusChecker:      nil, // Will be set via SetStatusChecker if needed
	}
}

// SetStatusChecker sets the gateway status checker
func (m *RoomServerAssignmentManager) SetStatusChecker(checker GatewayStatusChecker) {
	m.statusChecker = checker
}

// GetOrAssignServer returns the assigned server for a room or assigns a new one
// countryCode is used for regional routing (only for first assignment)
func (m *RoomServerAssignmentManager) GetOrAssignServer(ctx context.Context, roomID int64, countryCode string) (string, error) {
	// First, check if there's already an assignment
	var existingGatewayName string
	err := queries.Raw(
		"SELECT gateway_name FROM room_server_assignments WHERE room_id = $1",
		roomID,
	).QueryRow(m.db).Scan(&existingGatewayName)

	if err == nil {
		// Assignment exists, update last_used_at (sticky routing - ignore countryCode)
		_, err = queries.Raw(
			"UPDATE room_server_assignments SET last_used_at = $1 WHERE room_id = $2",
			time.Now().UTC(), roomID,
		).Exec(m.db)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to update last_used_at")
		}
		return existingGatewayName, nil
	}

	if err != sql.ErrNoRows {
		return "", pkgerr.Wrap(err, "query room assignment")
	}

	// No assignment exists - this is the first user in the room
	// Try to find regional servers first
	preferredServers := m.getPreferredServers(countryCode)

	serverLoads, err := m.getServerLoads(ctx)
	if err != nil {
		return "", pkgerr.Wrap(err, "get server loads")
	}

	reservedServers := m.getReservedServers(countryCode)
	selectedServer := m.selectLeastLoadedServer(serverLoads, preferredServers, reservedServers)

	// Create assignment with region info
	// Use INSERT ... ON CONFLICT to handle race conditions atomically
	var gatewayName string
	err = queries.Raw(`
		INSERT INTO room_server_assignments (room_id, gateway_name, region, assigned_at, last_used_at) 
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (room_id) 
		DO UPDATE SET last_used_at = EXCLUDED.last_used_at
		RETURNING gateway_name
	`, roomID, selectedServer, countryCode, time.Now().UTC(), time.Now().UTC()).QueryRow(m.db).Scan(&gatewayName)

	if err != nil {
		return "", pkgerr.Wrap(err, "insert room assignment")
	}
	
	// If gateway_name differs from selectedServer, it means another request won the race
	// Log this for debugging
	if gatewayName != selectedServer {
		log.Ctx(ctx).Debug().
			Int64("room_id", roomID).
			Str("selected_server", selectedServer).
			Str("actual_server", gatewayName).
			Msg("Race condition detected - using existing assignment")
	}

	log.Ctx(ctx).Info().
		Int64("room_id", roomID).
		Str("gateway_name", gatewayName).
		Str("country_code", countryCode).
		Bool("regional_match", len(preferredServers) > 0).
		Msg("Assigned room to server")

	return gatewayName, nil
}

// getPreferredServers returns list of servers for the given country code
func (m *RoomServerAssignmentManager) getPreferredServers(countryCode string) []string {
	if countryCode == "" {
		return nil
	}
	
	if servers, ok := m.serverRegions[countryCode]; ok {
		return servers
	}
	
	return nil
}

// getReservedServers returns list of servers reserved for OTHER regions
func (m *RoomServerAssignmentManager) getReservedServers(countryCode string) map[string]bool {
	reserved := make(map[string]bool)
	
	for region, servers := range m.serverRegions {
		if region == countryCode {
			continue // Skip own region
		}
		for _, server := range servers {
			reserved[server] = true
		}
	}
	
	return reserved
}

// getServerLoads calculates estimated load for each available server
// Load is calculated as: number_of_rooms * avgRoomOccupancy
func (m *RoomServerAssignmentManager) getServerLoads(ctx context.Context) (map[string]int, error) {
	loads := make(map[string]int)

	// Initialize all servers with 0 load
	for _, server := range m.availableServers {
		loads[server] = 0
	}

	// Count assigned rooms per server
	rows, err := queries.Raw(`
		SELECT gateway_name, COUNT(*) as room_count
		FROM room_server_assignments
		WHERE gateway_name = ANY($1)
		GROUP BY gateway_name
	`, "{"+strings.Join(m.availableServers, ",")+"}").Query(m.db)

	if err != nil {
		return nil, pkgerr.Wrap(err, "query server loads")
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var roomCount int
		if err := rows.Scan(&name, &roomCount); err != nil {
			return nil, pkgerr.Wrap(err, "scan server load")
		}
		// Estimate load: rooms * average occupancy
		loads[name] = roomCount * m.avgRoomOccupancy
	}

	if err := rows.Err(); err != nil {
		return nil, pkgerr.Wrap(err, "rows error")
	}

	return loads, nil
}

// selectLeastLoadedServer picks the server to fill sequentially
// Strategy: fill servers one by one to maxServerCapacity, then move to next
// If preferredServers is provided, it will try to select from those first
// reservedServers are excluded from selection (unless all non-reserved are at capacity)
func (m *RoomServerAssignmentManager) selectLeastLoadedServer(loads map[string]int, preferredServers []string, reservedServers map[string]bool) string {
	var selectedServer string
	maxLoad := -1

	// First, try to select from preferred servers (regional)
	if len(preferredServers) > 0 {
		for _, server := range preferredServers {
			// Check if server is in available list
			found := false
			for _, availServer := range m.availableServers {
				if availServer == server {
					found = true
					break
				}
			}
			if !found {
				continue
			}

			// Skip offline servers
			if m.statusChecker != nil && !m.statusChecker.IsGatewayOnline(server) {
				continue
			}

			load := loads[server]
			// Check if server has capacity for one more room
			if load+m.avgRoomOccupancy > m.maxServerCapacity {
				continue
			}
			
			// Select server with MAXIMUM load (to fill it first)
			if maxLoad == -1 || load > maxLoad {
				maxLoad = load
				selectedServer = server
			}
		}
		
		// If we found a regional server, use it
		if selectedServer != "" {
			return selectedServer
		}
	}

	// No preferred servers or all are at capacity - select from non-reserved servers
	maxLoad = -1
	for _, server := range m.availableServers {
		// Skip servers reserved for other regions
		if reservedServers[server] {
			continue
		}
		
		// Skip offline servers
		if m.statusChecker != nil && !m.statusChecker.IsGatewayOnline(server) {
			continue
		}
		
		load := loads[server]
		// Check if server has capacity for one more room
		if load+m.avgRoomOccupancy > m.maxServerCapacity {
			continue
		}
		
		// Select server with MAXIMUM load (to fill it first)
		if maxLoad == -1 || load > maxLoad {
			maxLoad = load
			selectedServer = server
		}
	}
	
	// If we found a non-reserved server, use it
	if selectedServer != "" {
		return selectedServer
	}

	// All non-reserved servers are at capacity
	// Final fallback: select least loaded non-reserved server (ignoring capacity check)
	minLoad := -1
	for _, server := range m.availableServers {
		// Still skip servers reserved for other regions
		if reservedServers[server] {
			continue
		}
		
		load := loads[server]
		if minLoad == -1 || load < minLoad {
			minLoad = load
			selectedServer = server
		}
	}

	return selectedServer
}

// UpdateLastUsed updates the last_used_at timestamp for a room assignment
func (m *RoomServerAssignmentManager) UpdateLastUsed(ctx context.Context, roomID int64) error {
	_, err := queries.Raw(
		"UPDATE room_server_assignments SET last_used_at = $1 WHERE room_id = $2",
		time.Now().UTC(), roomID,
	).Exec(m.db)

	if err != nil {
		return pkgerr.Wrap(err, "update last_used_at")
	}

	return nil
}

// CleanInactiveAssignments removes assignments for rooms with no active sessions
// Uses the same "active session" criteria as PeriodicSessionCleaner:
//  - removed_at IS NULL
//  - updated_at >= NOW() - DeadSessionPeriod
func (m *RoomServerAssignmentManager) CleanInactiveAssignments(ctx context.Context) error {
	// Calculate the cutoff time using the same logic as PeriodicSessionCleaner
	cutoffTime := time.Now().Add(-common.Config.DeadSessionPeriod)
	
	// Delete assignments where there are no ACTIVE sessions in the room
	// Active = removed_at IS NULL AND updated_at >= cutoff
	res, err := queries.Raw(`
		DELETE FROM room_server_assignments rsa
		WHERE NOT EXISTS (
			SELECT 1 FROM sessions s 
			WHERE s.room_id = rsa.room_id 
			AND s.removed_at IS NULL
			AND s.updated_at >= $1
		)
	`, cutoffTime).Exec(m.db)

	if err != nil {
		return pkgerr.Wrap(err, "clean inactive assignments")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return pkgerr.Wrap(err, "get rows affected")
	}

	if rowsAffected > 0 {
		log.Ctx(ctx).Info().
			Int64("cleaned_assignments", rowsAffected).
			Msg("Cleaned inactive room server assignments")
	}

	return nil
}

// MigrateServerAssignments moves all room assignments from one server to another
// Used for failover - migrates all rooms from failed server to failover server
func (m *RoomServerAssignmentManager) MigrateServerAssignments(ctx context.Context, fromServer, toServer string) (int, error) {
	res, err := queries.Raw(`
		UPDATE room_server_assignments
		SET gateway_name = $1, last_used_at = $2
		WHERE gateway_name = $3
	`, toServer, time.Now().UTC(), fromServer).Exec(m.db)
	
	if err != nil {
		return 0, pkgerr.Wrap(err, "migrate server assignments")
	}
	
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, pkgerr.Wrap(err, "get rows affected")
	}
	
	log.Ctx(ctx).Info().
		Str("from_server", fromServer).
		Str("to_server", toServer).
		Int64("count", rowsAffected).
		Msg("Migrated room server assignments")
	
	return int(rowsAffected), nil
}

// DistributeServerAssignments distributes room assignments from failed server among alive servers
// Used when failover servers are exhausted - distributes to alive primary servers
func (m *RoomServerAssignmentManager) DistributeServerAssignments(ctx context.Context, failedServer string, aliveServers []string) (int, error) {
	if len(aliveServers) == 0 {
		return 0, pkgerr.New("no alive servers provided")
	}
	
	// Get all assignments for failed server
	type Assignment struct {
		RoomID int64 `boil:"room_id"`
	}
	
	var assignments []Assignment
	err := queries.Raw(`
		SELECT room_id
		FROM room_server_assignments
		WHERE gateway_name = $1
		ORDER BY room_id
	`, failedServer).Bind(ctx, m.db, &assignments)
	
	if err != nil {
		return 0, pkgerr.Wrap(err, "get assignments")
	}
	
	if len(assignments) == 0 {
		return 0, nil
	}
	
	// Distribute assignments round-robin among alive servers
	now := time.Now().UTC()
	for i, assignment := range assignments {
		targetServer := aliveServers[i%len(aliveServers)]
		
		_, err := queries.Raw(`
			UPDATE room_server_assignments
			SET gateway_name = $1, last_used_at = $2
			WHERE room_id = $3
		`, targetServer, now, assignment.RoomID).Exec(m.db)
		
		if err != nil {
			log.Ctx(ctx).Error().
				Err(err).
				Int64("room_id", assignment.RoomID).
				Str("target_server", targetServer).
				Msg("Failed to distribute assignment")
			continue
		}
	}
	
	log.Ctx(ctx).Warn().
		Str("failed_server", failedServer).
		Strs("alive_servers", aliveServers).
		Int("count", len(assignments)).
		Msg("Distributed room server assignments (emergency mode)")
	
	return len(assignments), nil
}
