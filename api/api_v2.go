package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	pkgerr "github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/Bnei-Baruch/gxydb-api/common"
	"github.com/Bnei-Baruch/gxydb-api/middleware"
	"github.com/Bnei-Baruch/gxydb-api/models"
	"github.com/Bnei-Baruch/gxydb-api/pkg/httputil"
	"github.com/Bnei-Baruch/gxydb-api/pkg/sqlutil"
)

// webinarLanguageRe validates the language token used to build webinar room
// names. Restricting it to lowercase letters keeps room names predictable and
// makes it safe to interpolate into the room-name regex used in SQL.
var webinarLanguageRe = regexp.MustCompile(`^[a-z]+$`)

func (a *App) V2GetConfig(w http.ResponseWriter, r *http.Request) {
	cfg := V2Config{
		Gateways:      make(map[string]map[string]*V2Gateway),
		IceServers:    common.Config.IceServers,
		DynamicConfig: make(map[string]string),
	}

	gateways := a.cache.gateways.Values()
	for _, gateway := range gateways {
		if gateway.Disabled || gateway.RemovedAt.Valid {
			continue
		}

		respGateway := &V2Gateway{
			Name: gateway.Name,
			URL:  gateway.URL,
			Type: gateway.Type,
		}

		if cfg.Gateways[gateway.Type] == nil {
			cfg.Gateways[gateway.Type] = make(map[string]*V2Gateway)
		}
		cfg.Gateways[gateway.Type][gateway.Name] = respGateway
	}

	kvs := a.cache.dynamicConfig.Values()
	for _, kv := range kvs {
		cfg.DynamicConfig[kv.Key] = kv.Value
	}
	cfg.LastModified = a.cache.dynamicConfig.LastModified()

	httputil.RespondWithJSON(w, http.StatusOK, cfg)
}

func (a *App) V2GetRoomsStatistics(w http.ResponseWriter, r *http.Request) {
	stats, err := a.roomsStatisticsManager.GetAll()
	if err != nil {
		httputil.NewInternalError(err).Abort(w, r)
		return
	}

	data := make(map[string]*V2RoomStatistics, len(stats))
	for _, roomStats := range stats {
		// roomStats.RoomID now contains Janus room ID (gateway_uid) as string
		data[roomStats.RoomID] = &V2RoomStatistics{OnAir: roomStats.OnAir}
	}

	httputil.RespondWithJSON(w, http.StatusOK, data)
}

func (a *App) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()

	err := a.DB.(*sql.DB).PingContext(ctx)
	if err != nil {
		httputil.RespondWithError(w, http.StatusFailedDependency, fmt.Sprintf("DB ping: %s", err.Error()))
		return
	}

	if ctx.Err() == context.DeadlineExceeded {
		httputil.RespondWithError(w, http.StatusServiceUnavailable, "timeout")
		return
	}

	httputil.RespondSuccess(w)
}

type VHInfo struct {
	UserID    string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

func (a *App) V2GetVHInfo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	rCtx, ok := middleware.ContextFromRequest(r)
	if !ok {
		httputil.NewInternalError(errors.New("request missing context")).Abort(w, r)
		return
	}

	url := fmt.Sprintf("%s/profile/v1/profile/%s/short", common.Config.VHUrl, rCtx.IDClaims.Sub)
	payload, err := httputil.HTTPGetWithAuth(ctx, url, r.Header.Get("Authorization"))
	if err != nil {
		httputil.NewInternalError(err).Abort(w, r)
		return
	}

	var vhinfo VHInfo
	if err := json.Unmarshal([]byte(payload), &vhinfo); err != nil {
		httputil.NewInternalError(fmt.Errorf("json.Unmarshal VH info: %w", err)).Abort(w, r)
		return
	}

	httputil.RespondWithJSON(w, http.StatusOK, vhinfo)
}

func (a *App) V2GetRoomServer(w http.ResponseWriter, r *http.Request) {
	var req V2RoomServerRequest
	if err := httputil.DecodeJSONBody(w, r, &req); err != nil {
		err.Abort(w, r)
		return
	}

	// Validate room exists
	room, ok := a.cache.rooms.ByGatewayUID(req.Room)
	if !ok {
		httputil.NewNotFoundError().Abort(w, r)
		return
	}

	var gatewayName string
	var err error

	// Pinned rooms (SERVER_ROOMS): if the room is in the static pinned list,
	// the configured server always wins, bypassing both legacy and scale-mode
	// selection logic. The assignment is still recorded in room_server_assignments
	// so the rest of the system (rooms listing, load accounting) stays in sync.
	if pinned, ok := common.Config.ServerRooms[room.GatewayUID]; ok && pinned != "" {
		gatewayName, err = a.roomServerAssignmentManager.AssignPinnedServer(r.Context(), room.GatewayUID, pinned)
		if err != nil {
			httputil.NewInternalError(pkgerr.WithStack(err)).Abort(w, r)
			return
		}
		log.Ctx(r.Context()).Info().
			Str("room", req.Room).
			Str("janus", gatewayName).
			Msg("V2GetRoomServer pinned response (SERVER_ROOMS)")
		httputil.RespondWithJSON(w, http.StatusOK, V2RoomServerResponse{
			Janus: gatewayName,
		})
		return
	}

	// Check mode: legacy (use default gateway) or scale (load balancing)
	if !common.Config.ScaleMode {
		// Legacy mode: return room's default gateway
		gateway, ok := a.cache.gateways.ByID(room.DefaultGatewayID)
		if !ok {
			log.Ctx(r.Context()).Error().
				Int64("default_gateway_id", room.DefaultGatewayID).
				Int64("room_id", room.ID).
				Str("room_gateway_uid", req.Room).
				Msg("Gateway not found for room in legacy mode (SCALE=false). Ensure all rooms have valid default_gateway_id.")
			httputil.NewInternalError(pkgerr.Errorf("gateway not found for room %s", req.Room)).Abort(w, r)
			return
		}
		gatewayName = gateway.Name
	} else {
		// Scale mode: use load balancing with optional regional routing
		countryCode := ""
		if req.Geo != nil {
			countryCode = req.Geo.CountryCode
		}

		// Use Janus room ID (gateway_uid), not internal rooms.id
		janusRoomID := room.GatewayUID // Already string after migration
		gatewayName, err = a.roomServerAssignmentManager.GetOrAssignServer(r.Context(), janusRoomID, countryCode)
		if err != nil {
			httputil.NewInternalError(pkgerr.WithStack(err)).Abort(w, r)
			return
		}
	}

	log.Ctx(r.Context()).Info().
		Str("room", req.Room).
		Str("janus", gatewayName).
		Msg("V2GetRoomServer response")

	httputil.RespondWithJSON(w, http.StatusOK, V2RoomServerResponse{
		Janus: gatewayName,
	})
}

// V2GetWebinarRoomServer is the webinar mode counterpart of V2GetRoomServer.
// Instead of the client choosing the room, the client sends its language and the
// backend finds (or auto-creates) a room "<language>-<n>" with free capacity,
// then assigns a server using the same load balancing as galaxy scale mode.
func (a *App) V2GetWebinarRoomServer(w http.ResponseWriter, r *http.Request) {
	if common.Config.Mode != common.ModeWebinar {
		httputil.NewBadRequestError(nil, "webinar mode is not enabled").Abort(w, r)
		return
	}

	var req V2WebinarRoomServerRequest
	if err := httputil.DecodeJSONBody(w, r, &req); err != nil {
		err.Abort(w, r)
		return
	}

	language := strings.ToLower(strings.TrimSpace(req.Language))
	if !webinarLanguageRe.MatchString(language) {
		httputil.NewBadRequestError(nil, "invalid or missing language").Abort(w, r)
		return
	}

	countryCode := ""
	if req.Geo != nil {
		countryCode = req.Geo.CountryCode
	}

	// rooms.default_gateway_id is NOT NULL. In webinar mode the effective server
	// is taken from room_server_assignments (load balanced), so this value is just
	// a schema filler - pick the first available gateway.
	defaultGatewayID, ok := a.firstAvailableGatewayID()
	if !ok {
		httputil.NewInternalError(pkgerr.New("no available janus gateways configured")).Abort(w, r)
		return
	}

	var gatewayUID, roomName string
	var isNew bool

	// Find-or-create is serialized per language with a transaction scoped advisory
	// lock, so concurrent requests for the same language don't create duplicate rooms.
	//
	// NOTE: room occupancy is measured from active sessions in the DB. A session is
	// only created once the client actually joins the room (via MQTT events), so it
	// lags behind this assignment. Under a burst of simultaneous requests for the
	// same language, several users may therefore be packed into the same room beyond
	// WEBINAR_USERS_COUNT. This is an accepted trade-off of DB-based counting.
	err := sqlutil.InTx(r.Context(), a.DB, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(r.Context(),
			"SELECT pg_advisory_xact_lock(hashtext($1))", "webinar:"+language); err != nil {
			return pkgerr.Wrap(err, "advisory lock")
		}

		rows, err := tx.QueryContext(r.Context(), `
			SELECT r.gateway_uid, r.name,
			       (SELECT count(*) FROM sessions s WHERE s.room_id = r.gateway_uid AND s.removed_at IS NULL) AS cnt
			FROM rooms r
			WHERE r.disabled = false AND r.removed_at IS NULL AND r.name ~ ('^' || $1 || '-[0-9]+$')
		`, language)
		if err != nil {
			return pkgerr.Wrap(err, "query webinar rooms")
		}
		defer rows.Close()

		type webinarRoom struct {
			uid   string
			index int
			count int
		}
		var candidates []webinarRoom
		maxIndex := 0
		for rows.Next() {
			var uid, name string
			var cnt int
			if err := rows.Scan(&uid, &name, &cnt); err != nil {
				return pkgerr.Wrap(err, "scan webinar room")
			}
			idx := webinarRoomIndex(name, language)
			if idx <= 0 {
				continue
			}
			if idx > maxIndex {
				maxIndex = idx
			}
			candidates = append(candidates, webinarRoom{uid: uid, index: idx, count: cnt})
		}
		if err := rows.Err(); err != nil {
			return pkgerr.Wrap(err, "iterate webinar rooms")
		}

		// reuse the lowest-index room still under capacity
		sort.Slice(candidates, func(i, j int) bool { return candidates[i].index < candidates[j].index })
		for _, c := range candidates {
			if c.count < common.Config.WebinarUsersCount {
				gatewayUID = c.uid
				return nil
			}
		}

		// no room with free capacity - create the next one: <language>-<maxIndex+1>.
		// The Janus room ID (gateway_uid) IS the room name, e.g. "hebrew-1" (string
		// room IDs are supported after the janus_string_room_id migration). The
		// per-language advisory lock above guarantees the name is unique, so no
		// separate uid allocation / global lock is needed.
		roomName = fmt.Sprintf("%s-%d", language, maxIndex+1)

		room := models.Room{
			Name:             roomName,
			DefaultGatewayID: defaultGatewayID,
			GatewayUID:       roomName,
			Disabled:         false,
		}
		if err := room.Insert(tx, boil.Whitelist("name", "default_gateway_id", "gateway_uid", "disabled")); err != nil {
			return pkgerr.Wrap(err, "insert webinar room")
		}

		gatewayUID = roomName
		isNew = true
		return nil
	})
	if err != nil {
		httputil.NewInternalError(pkgerr.WithStack(err)).Abort(w, r)
		return
	}

	if isNew {
		// create the room on all gateways and refresh the cache so the rest of the
		// system (rooms listing, statistics, room_server) can see it.
		a.createVideoRoomOnGateways(common.Config.AvailableJanusServers, gatewayUID, roomName)
		if err := a.cache.rooms.Reload(a.DB); err != nil {
			log.Error().Err(err).Msg("Reload rooms cache after webinar room create")
		}
	}

	// assign (or stickily reuse) a server with the same load balancing as galaxy scale mode
	gatewayName, err := a.roomServerAssignmentManager.GetOrAssignServer(r.Context(), gatewayUID, countryCode)
	if err != nil {
		httputil.NewInternalError(pkgerr.WithStack(err)).Abort(w, r)
		return
	}

	log.Ctx(r.Context()).Info().
		Str("language", language).
		Str("room", gatewayUID).
		Str("janus", gatewayName).
		Bool("created", isNew).
		Msg("V2GetWebinarRoomServer response")

	httputil.RespondWithJSON(w, http.StatusOK, V2WebinarRoomServerResponse{
		Janus: gatewayName,
		Room:  gatewayUID,
	})
}

// webinarRoomIndex extracts the numeric suffix N from a room named
// "<language>-N". Returns 0 if the name doesn't match the expected shape.
func webinarRoomIndex(name, language string) int {
	prefix := language + "-"
	if !strings.HasPrefix(name, prefix) {
		return 0
	}
	idx, err := strconv.Atoi(name[len(prefix):])
	if err != nil {
		return 0
	}
	return idx
}

// firstAvailableGatewayID returns a gateway ID to use as the NOT NULL filler for
// rooms.default_gateway_id (in webinar mode the effective server comes from
// room_server_assignments, so the exact value here is not important).
// It prefers a gateway listed in AVAILABLE_JANUS_SERVERS, and otherwise falls back
// to any enabled rooms-type gateway in the cache.
func (a *App) firstAvailableGatewayID() (int64, bool) {
	for _, name := range common.Config.AvailableJanusServers {
		if g, ok := a.cache.gateways.ByName(name); ok && !g.Disabled && !g.RemovedAt.Valid {
			return g.ID, true
		}
	}
	for _, g := range a.cache.gateways.Values() {
		if g.Disabled || g.RemovedAt.Valid || g.Type != common.GatewayTypeRooms {
			continue
		}
		return g.ID, true
	}
	return 0, false
}
