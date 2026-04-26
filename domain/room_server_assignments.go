package domain

import (
	"context"
	"database/sql"
	"strings"
	"time"

	pkgerr "github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"

	"github.com/Bnei-Baruch/gxydb-api/common"
)

type RoomServerAssignment struct {
	RoomID      string
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
	statusChecker      GatewayStatusChecker
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

// SetStatusChecker sets the gateway status checker
func (m *RoomServerAssignmentManager) SetStatusChecker(checker GatewayStatusChecker) {
	m.statusChecker = checker
}

// GetOrAssignServer returns the assigned server for a room or assigns a new one
// countryCode is used for regional routing (only for first assignment)
func (m *RoomServerAssignmentManager) GetOrAssignServer(ctx context.Context, roomID string, countryCode string) (string, error) {
	// First, check if there's already an assignment
	var existingGatewayName string
	err := queries.Raw(
		"SELECT gateway_name FROM room_server_assignments WHERE room_id = $1",
		roomID,
	).QueryRow(m.db).Scan(&existingGatewayName)

	if err == nil {
		// Assignment exists - check if the assigned server is still online
		if m.statusChecker != nil && !m.statusChecker.IsGatewayOnline(existingGatewayName) {
			log.Ctx(ctx).Warn().
				Str("room_id", roomID).
				Str("offline_server", existingGatewayName).
				Msg("Assigned server is offline - reassigning room")

			serverLoads, loadErr := m.getServerLoads(ctx)
			if loadErr != nil {
				log.Ctx(ctx).Error().Err(loadErr).Msg("Failed to get server loads for reassignment")
				return existingGatewayName, nil
			}

			preferredServers := m.getPreferredServers(countryCode)
			reservedServers := m.getReservedServers(countryCode)
			newServer := m.selectLeastLoadedServer(serverLoads, preferredServers, reservedServers)

			if newServer != "" && newServer != existingGatewayName {
				if result, err := m.reassignRoom(ctx, roomID, existingGatewayName, newServer); result != "" {
					return result, err
				}
			} else {
				log.Ctx(ctx).Error().
					Str("room_id", roomID).
					Str("offline_server", existingGatewayName).
					Msg("No online servers available for reassignment")
			}
		}

		// Server is online (or reassignment failed) - sticky routing
		log.Ctx(ctx).Info().
			Str("room_id", roomID).
			Str("gateway", existingGatewayName).
			Msg("Sticky routing - returning existing assignment")
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
			Str("room_id", roomID).
			Str("selected_server", selectedServer).
			Str("actual_server", gatewayName).
			Msg("Race condition detected - using existing assignment")
	}

	log.Ctx(ctx).Info().
		Str("room_id", roomID).
		Str("gateway_name", gatewayName).
		Str("country_code", countryCode).
		Bool("regional_match", len(preferredServers) > 0).
		Msg("Assigned room to server")

	return gatewayName, nil
}

// AssignPinnedServer upserts a room->server assignment for a statically pinned room
// (configured via SERVER_ROOMS). Unlike GetOrAssignServer, it doesn't perform any
// load/online checks - the configured server always wins. If a different assignment
// already exists for the room, gateway_name is overwritten so the table stays in sync
// with the current SERVER_ROOMS configuration.
func (m *RoomServerAssignmentManager) AssignPinnedServer(ctx context.Context, roomID, server string) (string, error) {
	now := time.Now().UTC()

	var gatewayName string
	err := queries.Raw(`
		INSERT INTO room_server_assignments (room_id, gateway_name, region, assigned_at, last_used_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (room_id)
		DO UPDATE SET gateway_name = EXCLUDED.gateway_name, last_used_at = EXCLUDED.last_used_at
		RETURNING gateway_name
	`, roomID, server, "pinned", now, now).QueryRow(m.db).Scan(&gatewayName)

	if err != nil {
		return "", pkgerr.Wrap(err, "upsert pinned room assignment")
	}

	log.Ctx(ctx).Info().
		Str("room_id", roomID).
		Str("gateway_name", gatewayName).
		Msg("Pinned room assigned to server (SERVER_ROOMS)")

	return gatewayName, nil
}

// reassignRoom atomically reassigns a room from oldServer to newServer.
// Uses a transaction so pgpool routes all queries (including SELECTs) to the primary node,
// avoiding stale reads from replicas due to replication lag.
func (m *RoomServerAssignmentManager) reassignRoom(ctx context.Context, roomID, oldServer, newServer string) (string, error) {
	sqlDB, ok := m.db.(*sql.DB)
	if !ok {
		if executor, ok2 := m.db.(boil.Beginner); ok2 {
			return m.reassignRoomWithBeginner(ctx, roomID, oldServer, newServer, executor)
		}
		log.Ctx(ctx).Error().Msg("DB does not support transactions - falling back to non-transactional reassign")
		return m.reassignRoomDirect(ctx, roomID, oldServer, newServer)
	}

	tx, txErr := sqlDB.BeginTx(ctx, nil)
	if txErr != nil {
		log.Ctx(ctx).Error().Err(txErr).Msg("Failed to begin transaction for reassignment")
		return m.reassignRoomDirect(ctx, roomID, oldServer, newServer)
	}
	defer tx.Rollback()

	now := time.Now().UTC()

	// All queries in this tx go to primary (pgpool routes tx with writes to primary)
	result, err := tx.ExecContext(ctx,
		"UPDATE room_server_assignments SET gateway_name = $1, last_used_at = $2 WHERE room_id = $3 AND gateway_name = $4",
		newServer, now, roomID, oldServer)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to update room reassignment")
		return "", nil
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected > 0 {
		if commitErr := tx.Commit(); commitErr != nil {
			log.Ctx(ctx).Error().Err(commitErr).Msg("Failed to commit reassignment")
			return "", nil
		}
		log.Ctx(ctx).Info().
			Str("room_id", roomID).
			Str("old_server", oldServer).
			Str("new_server", newServer).
			Msg("Room reassigned from offline server")
		return newServer, nil
	}

	// Another request beat us — read current value FROM PRIMARY (same tx)
	var currentGateway string
	if err := tx.QueryRowContext(ctx,
		"SELECT gateway_name FROM room_server_assignments WHERE room_id = $1",
		roomID).Scan(&currentGateway); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to read current assignment in tx")
		return "", nil
	}
	tx.Commit()

	if m.statusChecker == nil || m.statusChecker.IsGatewayOnline(currentGateway) {
		log.Ctx(ctx).Info().
			Str("room_id", roomID).
			Str("current_server", currentGateway).
			Msg("Room already reassigned by concurrent request")
		return currentGateway, nil
	}

	// Current server is still offline — force reassign in a new tx
	log.Ctx(ctx).Warn().
		Str("room_id", roomID).
		Str("current_server", currentGateway).
		Str("new_server", newServer).
		Msg("Concurrent reassignment still points to offline server - force reassign")
	return m.reassignRoom(ctx, roomID, currentGateway, newServer)
}

func (m *RoomServerAssignmentManager) reassignRoomWithBeginner(ctx context.Context, roomID, oldServer, newServer string, beginner boil.Beginner) (string, error) {
	tx, err := beginner.Begin()
	if err != nil {
		return m.reassignRoomDirect(ctx, roomID, oldServer, newServer)
	}
	defer tx.Rollback()

	now := time.Now().UTC()
	result, err := queries.Raw(
		"UPDATE room_server_assignments SET gateway_name = $1, last_used_at = $2 WHERE room_id = $3 AND gateway_name = $4",
		newServer, now, roomID, oldServer).Exec(tx)
	if err != nil {
		return "", nil
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected > 0 {
		tx.Commit()
		log.Ctx(ctx).Info().
			Str("room_id", roomID).Str("old_server", oldServer).Str("new_server", newServer).
			Msg("Room reassigned from offline server")
		return newServer, nil
	}

	var currentGateway string
	if err := queries.Raw(
		"SELECT gateway_name FROM room_server_assignments WHERE room_id = $1",
		roomID).QueryRow(tx).Scan(&currentGateway); err != nil {
		return "", nil
	}
	tx.Commit()

	if m.statusChecker == nil || m.statusChecker.IsGatewayOnline(currentGateway) {
		log.Ctx(ctx).Info().
			Str("room_id", roomID).Str("current_server", currentGateway).
			Msg("Room already reassigned by concurrent request")
		return currentGateway, nil
	}

	log.Ctx(ctx).Warn().
		Str("room_id", roomID).Str("current_server", currentGateway).Str("new_server", newServer).
		Msg("Concurrent reassignment still points to offline server - force reassign")
	return m.reassignRoom(ctx, roomID, currentGateway, newServer)
}

// reassignRoomDirect is a fallback without transactions
func (m *RoomServerAssignmentManager) reassignRoomDirect(ctx context.Context, roomID, oldServer, newServer string) (string, error) {
	now := time.Now().UTC()
	result, err := queries.Raw(
		"UPDATE room_server_assignments SET gateway_name = $1, last_used_at = $2 WHERE room_id = $3 AND gateway_name = $4",
		newServer, now, roomID, oldServer).Exec(m.db)
	if err != nil {
		return "", nil
	}
	if rowsAffected, _ := result.RowsAffected(); rowsAffected > 0 {
		log.Ctx(ctx).Info().
			Str("room_id", roomID).Str("old_server", oldServer).Str("new_server", newServer).
			Msg("Room reassigned from offline server")
		return newServer, nil
	}
	return "", nil
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

// selectLeastLoadedServer picks the online server with the least load.
// If preferredServers is provided, it will try to select from those first.
// reservedServers are excluded from selection (unless all non-reserved are at capacity).
func (m *RoomServerAssignmentManager) selectLeastLoadedServer(loads map[string]int, preferredServers []string, reservedServers map[string]bool) string {
	var selectedServer string
	minLoad := -1

	// First, try to select from preferred servers (regional)
	if len(preferredServers) > 0 {
		for _, server := range preferredServers {
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

			if m.statusChecker != nil && !m.statusChecker.IsGatewayOnline(server) {
				continue
			}

			load := loads[server]
			if load+m.avgRoomOccupancy > m.maxServerCapacity {
				continue
			}

			if minLoad == -1 || load < minLoad {
				minLoad = load
				selectedServer = server
			}
		}

		if selectedServer != "" {
			return selectedServer
		}
	}

	// No preferred servers or all are at capacity - select from non-reserved servers
	minLoad = -1
	for _, server := range m.availableServers {
		if reservedServers[server] {
			continue
		}

		if m.statusChecker != nil && !m.statusChecker.IsGatewayOnline(server) {
			continue
		}

		load := loads[server]
		if load+m.avgRoomOccupancy > m.maxServerCapacity {
			continue
		}

		if minLoad == -1 || load < minLoad {
			minLoad = load
			selectedServer = server
		}
	}

	if selectedServer != "" {
		return selectedServer
	}

	// All non-reserved servers are at capacity - fallback to least loaded online server
	minLoad = -1
	for _, server := range m.availableServers {
		if reservedServers[server] {
			continue
		}
		if m.statusChecker != nil && !m.statusChecker.IsGatewayOnline(server) {
			continue
		}

		load := loads[server]
		if minLoad == -1 || load < minLoad {
			minLoad = load
			selectedServer = server
		}
	}

	if selectedServer != "" {
		return selectedServer
	}

	// No online servers found (e.g. status data not yet available after restart).
	// Last resort: pick least loaded server ignoring online status.
	minLoad = -1
	for _, server := range m.availableServers {
		load := loads[server]
		if minLoad == -1 || load < minLoad {
			minLoad = load
			selectedServer = server
		}
	}

	return selectedServer
}

// UpdateLastUsed updates the last_used_at timestamp for a room assignment.
// Accepts an optional executor (e.g. *sql.Tx) to reuse the caller's transaction
// and avoid acquiring a second connection from the pool (which can deadlock).
func (m *RoomServerAssignmentManager) UpdateLastUsed(ctx context.Context, roomID string, exec ...boil.Executor) error {
	var e boil.Executor = m.db
	if len(exec) > 0 && exec[0] != nil {
		e = exec[0]
	}
	_, err := queries.Raw(
		"UPDATE room_server_assignments SET last_used_at = $1 WHERE room_id = $2",
		time.Now().UTC(), roomID,
	).Exec(e)

	if err != nil {
		return pkgerr.Wrap(err, "update last_used_at")
	}

	return nil
}

// CleanInactiveAssignments removes assignments for rooms with no active sessions
// Uses the same "active session" criteria as PeriodicSessionCleaner:
//   - removed_at IS NULL
//   - updated_at >= NOW() - DeadSessionPeriod
func (m *RoomServerAssignmentManager) CleanInactiveAssignments(ctx context.Context, exec ...boil.Executor) error {
	var e boil.Executor = m.db
	if len(exec) > 0 && exec[0] != nil {
		e = exec[0]
	}

	cutoffTime := time.Now().Add(-common.Config.DeadSessionPeriod)

	res, err := queries.Raw(`
		DELETE FROM room_server_assignments rsa
		WHERE NOT EXISTS (
			SELECT 1 FROM sessions s 
			WHERE s.room_id = rsa.room_id 
			AND s.removed_at IS NULL
			AND s.updated_at >= $1
		)
	`, cutoffTime).Exec(e)

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
		RoomID string `boil:"room_id"`
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
			WHERE room_id = $3 AND gateway_name = $4
		`, targetServer, now, assignment.RoomID, failedServer).Exec(m.db)

		if err != nil {
			log.Ctx(ctx).Error().
				Err(err).
				Str("room_id", assignment.RoomID).
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
