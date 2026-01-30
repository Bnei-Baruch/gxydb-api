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

type RoomServerAssignmentManager struct {
	db                 common.DBInterface
	availableServers   []string
	maxServerCapacity  int
	avgRoomOccupancy   int
	serverRegions      map[string][]string
}

func NewRoomServerAssignmentManager(db common.DBInterface, servers []string, maxCapacity, avgOccupancy int, regions map[string][]string) *RoomServerAssignmentManager {
	return &RoomServerAssignmentManager{
		db:                 db,
		availableServers:   servers,
		maxServerCapacity:  maxCapacity,
		avgRoomOccupancy:   avgOccupancy,
		serverRegions:      regions,
	}
}

// GetOrAssignServer returns the assigned server for a room or assigns a new one
// countryCode is used for regional routing (only for first assignment)
func (m *RoomServerAssignmentManager) GetOrAssignServer(ctx context.Context, roomID int64, countryCode string) (string, error) {
	// First, check if there's already an assignment
	var gatewayName string
	err := queries.Raw(
		"SELECT gateway_name FROM room_server_assignments WHERE room_id = $1",
		roomID,
	).QueryRow(m.db).Scan(&gatewayName)

	if err == nil {
		// Assignment exists, update last_used_at (sticky routing - ignore countryCode)
		_, err = queries.Raw(
			"UPDATE room_server_assignments SET last_used_at = $1 WHERE room_id = $2",
			time.Now().UTC(), roomID,
		).Exec(m.db)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to update last_used_at")
		}
		return gatewayName, nil
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

	selectedServer := m.selectLeastLoadedServer(serverLoads, preferredServers)

	// Create assignment with region info
	_, err = queries.Raw(
		"INSERT INTO room_server_assignments (room_id, gateway_name, region, assigned_at, last_used_at) VALUES ($1, $2, $3, $4, $5)",
		roomID, selectedServer, countryCode, time.Now().UTC(), time.Now().UTC(),
	).Exec(m.db)

	if err != nil {
		return "", pkgerr.Wrap(err, "insert room assignment")
	}

	log.Ctx(ctx).Info().
		Int64("room_id", roomID).
		Str("gateway_name", selectedServer).
		Str("country_code", countryCode).
		Bool("regional_match", len(preferredServers) > 0).
		Msg("Assigned room to server")

	return selectedServer, nil
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
func (m *RoomServerAssignmentManager) selectLeastLoadedServer(loads map[string]int, preferredServers []string) string {
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

	// No preferred servers or all are at capacity - select from all available
	maxLoad = -1
	for _, server := range m.availableServers {
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

	// If all servers are at capacity, select least loaded anyway (fallback)
	if selectedServer == "" {
		minLoad := -1
		for _, server := range m.availableServers {
			load := loads[server]
			if minLoad == -1 || load < minLoad {
				minLoad = load
				selectedServer = server
			}
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
func (m *RoomServerAssignmentManager) CleanInactiveAssignments(ctx context.Context) error {
	// Delete assignments where there are no active sessions in the room
	res, err := queries.Raw(`
		DELETE FROM room_server_assignments rsa
		WHERE NOT EXISTS (
			SELECT 1 FROM sessions s 
			WHERE s.room_id = rsa.room_id 
			AND s.removed_at IS NULL
		)
	`).Exec(m.db)

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
