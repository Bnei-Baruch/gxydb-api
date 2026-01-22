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
}

func NewRoomServerAssignmentManager(db common.DBInterface, servers []string, maxCapacity, avgOccupancy int) *RoomServerAssignmentManager {
	return &RoomServerAssignmentManager{
		db:                 db,
		availableServers:   servers,
		maxServerCapacity:  maxCapacity,
		avgRoomOccupancy:   avgOccupancy,
	}
}

// GetOrAssignServer returns the assigned server for a room or assigns a new one
func (m *RoomServerAssignmentManager) GetOrAssignServer(ctx context.Context, roomID int64) (string, error) {
	// First, check if there's already an assignment
	var gatewayName string
	err := queries.Raw(
		"SELECT gateway_name FROM room_server_assignments WHERE room_id = $1",
		roomID,
	).QueryRow(m.db).Scan(&gatewayName)

	if err == nil {
		// Assignment exists, update last_used_at
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

	// No assignment exists, find the least loaded server
	serverLoads, err := m.getServerLoads(ctx)
	if err != nil {
		return "", pkgerr.Wrap(err, "get server loads")
	}

	selectedServer := m.selectLeastLoadedServer(serverLoads)

	// Create assignment
	_, err = queries.Raw(
		"INSERT INTO room_server_assignments (room_id, gateway_name, assigned_at, last_used_at) VALUES ($1, $2, $3, $4)",
		roomID, selectedServer, time.Now().UTC(), time.Now().UTC(),
	).Exec(m.db)

	if err != nil {
		return "", pkgerr.Wrap(err, "insert room assignment")
	}

	log.Ctx(ctx).Info().
		Int64("room_id", roomID).
		Str("gateway_name", selectedServer).
		Msg("Assigned room to server")

	return selectedServer, nil
}

// getServerLoads calculates current load for each available server
func (m *RoomServerAssignmentManager) getServerLoads(ctx context.Context) (map[string]int, error) {
	loads := make(map[string]int)

	// Initialize all servers with 0 load
	for _, server := range m.availableServers {
		loads[server] = 0
	}

	// Count active sessions per gateway
	rows, err := queries.Raw(`
		SELECT g.name, COUNT(DISTINCT s.user_id) as load
		FROM gateways g
		LEFT JOIN sessions s ON s.gateway_id = g.id AND s.removed_at IS NULL
		WHERE g.name = ANY($1) AND g.disabled = false AND g.removed_at IS NULL
		GROUP BY g.name
	`, "{"+strings.Join(m.availableServers, ",")+"}").Query(m.db)

	if err != nil {
		return nil, pkgerr.Wrap(err, "query server loads")
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var load int
		if err := rows.Scan(&name, &load); err != nil {
			return nil, pkgerr.Wrap(err, "scan server load")
		}
		loads[name] = load
	}

	if err := rows.Err(); err != nil {
		return nil, pkgerr.Wrap(err, "rows error")
	}

	return loads, nil
}

// selectLeastLoadedServer picks the server with the lowest load
func (m *RoomServerAssignmentManager) selectLeastLoadedServer(loads map[string]int) string {
	var selectedServer string
	minLoad := -1

	for _, server := range m.availableServers {
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
