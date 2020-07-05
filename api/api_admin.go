package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	janus_admin "github.com/edoshor/janus-go/admin"
	janus_plugins "github.com/edoshor/janus-go/plugins"
	"github.com/gorilla/mux"
	pkgerr "github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/gxydb-api/common"
	"github.com/Bnei-Baruch/gxydb-api/domain"
	"github.com/Bnei-Baruch/gxydb-api/middleware"
	"github.com/Bnei-Baruch/gxydb-api/models"
	"github.com/Bnei-Baruch/gxydb-api/pkg/httputil"
	"github.com/Bnei-Baruch/gxydb-api/pkg/mathutil"
	"github.com/Bnei-Baruch/gxydb-api/pkg/sqlutil"
)

func (a *App) AdminGatewaysHandleInfo(w http.ResponseWriter, r *http.Request) {
	if !common.Config.SkipPermissions && !middleware.RequestHasRole(r, common.RoleAdmin, common.RoleRoot) {
		httputil.NewForbiddenError().Abort(w, r)
		return
	}

	vars := mux.Vars(r)
	gatewayID := vars["gateway_id"]
	gateway, ok := a.cache.gateways.ByName(gatewayID)
	if !ok {
		httputil.NewNotFoundError().Abort(w, r)
		return
	}

	sessionIDStr := vars["session_id"]
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 64)
	if err != nil {
		httputil.NewNotFoundError().Abort(w, r)
		return
	}

	handleIDStr := vars["handle_id"]
	handleID, err := strconv.ParseUint(handleIDStr, 10, 64)
	if err != nil {
		httputil.NewNotFoundError().Abort(w, r)
		return
	}

	api, err := domain.GatewayAdminAPIRegistry.For(gateway)
	if err != nil {
		httputil.NewInternalError(pkgerr.WithMessage(err, "init admin api")).Abort(w, r)
		return
	}

	info, err := api.HandleInfo(sessionID, handleID)
	if err != nil {
		var tErr *janus_admin.ErrorAMResponse
		if errors.As(err, &tErr) {
			if tErr.Err.Code == 458 || tErr.Err.Code == 459 { // no such session or no such handle
				httputil.NewNotFoundError().Abort(w, r)
				return
			}
		}
		httputil.NewInternalError(pkgerr.Wrap(err, "api.HandleInfo")).Abort(w, r)
		return
	}

	httputil.RespondWithJSON(w, http.StatusOK, info)
}

func (a *App) AdminListRooms(w http.ResponseWriter, r *http.Request) {
	if !common.Config.SkipPermissions && !middleware.RequestHasRole(r, common.RoleRoot) {
		httputil.NewForbiddenError().Abort(w, r)
		return
	}

	listParams, err := ParseListParams(r)
	if err != nil {
		httputil.NewBadRequestError(err, "malformed list parameters").Abort(w, r)
		return
	}

	mods := make([]qm.QueryMod, 0)

	// count query
	var total int64
	countMods := append([]qm.QueryMod{qm.Select("count(DISTINCT id)")}, mods...)
	err = models.Rooms(countMods...).QueryRow(a.DB).Scan(&total)
	if err != nil {
		httputil.NewInternalError(err).Abort(w, r)
		return
	} else if total == 0 {
		httputil.RespondWithJSON(w, http.StatusOK, RoomsResponse{Rooms: make([]*models.Room, 0)})
		return
	}

	// order, limit, offset
	_, offset := listParams.appendListMods(&mods)
	if int64(offset) >= total {
		httputil.RespondWithJSON(w, http.StatusOK, RoomsResponse{Rooms: make([]*models.Room, 0)})
		return
	}

	// data query
	rooms, err := models.Rooms(mods...).All(a.DB)
	if err != nil {
		httputil.NewInternalError(err).Abort(w, r)
		return
	}

	httputil.RespondWithJSON(w, http.StatusOK, RoomsResponse{
		ListResponse: ListResponse{
			Total: total,
		},
		Rooms: rooms,
	})
}

func (a *App) AdminCreateRoom(w http.ResponseWriter, r *http.Request) {
	if !common.Config.SkipPermissions && !middleware.RequestHasRole(r, common.RoleRoot) {
		httputil.NewForbiddenError().Abort(w, r)
		return
	}

	var data models.Room
	if err := httputil.DecodeJSONBody(w, r, &data); err != nil {
		err.Abort(w, r)
		return
	}
	a.requestContext(r).Params = data

	if data.GatewayUID <= 0 {
		httputil.NewBadRequestError(nil, "gateway_uid must be a positive integer").Abort(w, r)
		return
	}

	if _, ok := a.cache.gateways.ByID(data.DefaultGatewayID); !ok {
		httputil.NewBadRequestError(nil, "gateway doesn't exists").Abort(w, r)
		return
	}

	// TODO: gateway_uid should be fully managed by us and not user input !
	if _, ok := a.cache.rooms.ByGatewayUID(data.GatewayUID); ok {
		httputil.NewBadRequestError(nil, "room already exists [gateway_uid]").Abort(w, r)
		return
	}

	if exists, _ := models.Rooms(models.RoomWhere.Name.EQ(data.Name)).Exists(a.DB); exists {
		httputil.NewBadRequestError(nil, "room already exists [name]").Abort(w, r)
		return
	}

	err := sqlutil.InTx(r.Context(), a.DB, func(tx *sql.Tx) error {
		// create room in DB
		if err := data.Insert(a.DB, boil.Whitelist("name", "default_gateway_id", "gateway_uid", "disabled")); err != nil {
			return pkgerr.WithStack(err)
		}

		// create room in gateways
		room := &janus_plugins.VideoroomRoom{
			Room:               data.GatewayUID,
			Description:        data.Name,
			Secret:             common.Config.GatewayRoomsSecret,
			Publishers:         100,
			Bitrate:            64000,
			FirFreq:            10,
			AudioCodec:         "opus",
			VideoCodec:         "h264",
			AudioLevelExt:      true,
			AudioLevelEvent:    true,
			AudioActivePackets: 25,
			AudioLevelAverage:  100,
			VideoOrientExt:     true,
			PlayoutDelayExt:    true,
			TransportWideCCExt: true,
		}
		request := janus_plugins.MakeVideoroomRequestFactory(common.Config.GatewayVideoroomAdminKey).
			CreateRequest(room, true, []string{})

		for _, gateway := range a.cache.gateways.Values() {
			if gateway.Type != common.GatewayTypeRooms {
				continue
			}

			api, err := domain.GatewayAdminAPIRegistry.For(gateway)
			if err != nil {
				return pkgerr.WithMessage(err, "Admin API for gateway")
			}

			if _, err = api.MessagePlugin(request); err != nil {
				return pkgerr.Wrap(err, "api.MessagePlugin")
			}
		}

		return nil
	})

	if err != nil {
		var hErr *httputil.HttpError
		if errors.As(err, &hErr) {
			hErr.Abort(w, r)
		} else {
			httputil.NewInternalError(err).Abort(w, r)
		}
		return
	}

	if err := a.cache.rooms.Reload(a.DB); err != nil {
		log.Error().Err(err).Msg("Reload cache")
	}

	httputil.RespondWithJSON(w, http.StatusCreated, data)
}

func (a *App) AdminGetRoom(w http.ResponseWriter, r *http.Request) {
	if !common.Config.SkipPermissions && !middleware.RequestHasRole(r, common.RoleRoot) {
		httputil.NewForbiddenError().Abort(w, r)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		httputil.NewNotFoundError().Abort(w, r)
		return
	}

	room, err := models.FindRoom(a.DB, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httputil.NewNotFoundError().Abort(w, r)
		} else {
			httputil.NewInternalError(pkgerr.WithStack(err)).Abort(w, r)
		}
		return
	}

	httputil.RespondWithJSON(w, http.StatusOK, room)
}

func (a *App) AdminUpdateRoom(w http.ResponseWriter, r *http.Request) {
	if !common.Config.SkipPermissions && !middleware.RequestHasRole(r, common.RoleRoot) {
		httputil.NewForbiddenError().Abort(w, r)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		httputil.NewNotFoundError().Abort(w, r)
		return
	}

	room, err := models.FindRoom(a.DB, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httputil.NewNotFoundError().Abort(w, r)
		} else {
			httputil.NewInternalError(pkgerr.WithStack(err)).Abort(w, r)
		}
		return
	}

	var data models.Room
	if err := httputil.DecodeJSONBody(w, r, &data); err != nil {
		err.Abort(w, r)
		return
	}
	a.requestContext(r).Params = data

	if data.GatewayUID <= 0 {
		httputil.NewBadRequestError(nil, "gateway_uid must be a positive integer").Abort(w, r)
		return
	}

	if _, ok := a.cache.gateways.ByID(data.DefaultGatewayID); !ok {
		httputil.NewBadRequestError(nil, "gateway doesn't exists").Abort(w, r)
		return
	}

	if exists, _ := models.Rooms(models.RoomWhere.GatewayUID.EQ(data.GatewayUID), models.RoomWhere.ID.NEQ(room.ID)).Exists(a.DB); exists {
		httputil.NewBadRequestError(nil, "room already exists [gateway_uid]").Abort(w, r)
		return
	}

	if exists, _ := models.Rooms(models.RoomWhere.Name.EQ(data.Name), models.RoomWhere.ID.NEQ(room.ID)).Exists(a.DB); exists {
		httputil.NewBadRequestError(nil, "room already exists [name]").Abort(w, r)
		return
	}

	err = sqlutil.InTx(r.Context(), a.DB, func(tx *sql.Tx) error {
		shouldUpdateGateways := room.Name != data.Name &&
			!room.RemovedAt.Valid

		// update room in DB
		room.Name = data.Name
		room.DefaultGatewayID = data.DefaultGatewayID
		room.Disabled = data.Disabled
		room.UpdatedAt = null.TimeFrom(time.Now().UTC())
		if _, err := room.Update(a.DB, boil.Whitelist("name", "default_gateway_id", "disabled", "updated_at")); err != nil {
			return pkgerr.WithStack(err)
		}

		if !shouldUpdateGateways {
			return nil
		}

		// update room in gateways
		room := &janus_plugins.VideoroomRoomForEdit{
			Room:        data.GatewayUID,
			Description: data.Name,
		}
		request := janus_plugins.MakeVideoroomRequestFactory(common.Config.GatewayVideoroomAdminKey).
			EditRequest(room, true, common.Config.GatewayRoomsSecret)

		for _, gateway := range a.cache.gateways.Values() {
			if gateway.Type != common.GatewayTypeRooms {
				continue
			}

			api, err := domain.GatewayAdminAPIRegistry.For(gateway)
			if err != nil {
				return pkgerr.WithMessage(err, "Admin API for gateway")
			}

			if _, err = api.MessagePlugin(request); err != nil {
				return pkgerr.Wrap(err, "api.MessagePlugin")
			}
		}

		return nil
	})

	if err != nil {
		var hErr *httputil.HttpError
		if errors.As(err, &hErr) {
			hErr.Abort(w, r)
		} else {
			httputil.NewInternalError(err).Abort(w, r)
		}
		return
	}

	if err := a.cache.rooms.Reload(a.DB); err != nil {
		log.Error().Err(err).Msg("Reload cache")
	}

	httputil.RespondWithJSON(w, http.StatusOK, room)
}

func (a *App) AdminDeleteRoom(w http.ResponseWriter, r *http.Request) {
	if !common.Config.SkipPermissions && !middleware.RequestHasRole(r, common.RoleRoot) {
		httputil.NewForbiddenError().Abort(w, r)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		httputil.NewNotFoundError().Abort(w, r)
		return
	}

	room, err := models.FindRoom(a.DB, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httputil.NewNotFoundError().Abort(w, r)
		} else {
			httputil.NewInternalError(pkgerr.WithStack(err)).Abort(w, r)
		}
		return
	}

	err = sqlutil.InTx(r.Context(), a.DB, func(tx *sql.Tx) error {
		room.RemovedAt = null.TimeFrom(time.Now().UTC())
		if _, err := room.Update(a.DB, boil.Whitelist(models.RoomColumns.RemovedAt)); err != nil {
			return httputil.NewInternalError(pkgerr.WithStack(err))
		}

		request := janus_plugins.MakeVideoroomRequestFactory(common.Config.GatewayVideoroomAdminKey).
			DestroyRequest(room.GatewayUID, true, common.Config.GatewayRoomsSecret)

		for _, gateway := range a.cache.gateways.Values() {
			if gateway.Type != common.GatewayTypeRooms {
				continue
			}

			api, err := domain.GatewayAdminAPIRegistry.For(gateway)
			if err != nil {
				return pkgerr.WithMessage(err, "Admin API for gateway")
			}

			if _, err = api.MessagePlugin(request); err != nil {
				return pkgerr.Wrap(err, "api.MessagePlugin")
			}
		}

		return nil
	})

	if err != nil {
		var hErr *httputil.HttpError
		if errors.As(err, &hErr) {
			hErr.Abort(w, r)
		} else {
			httputil.NewInternalError(err).Abort(w, r)
		}
		return
	}

	if err := a.cache.rooms.Reload(a.DB); err != nil {
		log.Error().Err(err).Msg("Reload cache")
	}

	httputil.RespondSuccess(w)
}

type ListParams struct {
	PageNumber int    `json:"page_no"`
	PageSize   int    `json:"page_size"`
	OrderBy    string `json:"order_by"`
	GroupBy    string `json:"-"`
}

func (p *ListParams) appendListMods(mods *[]qm.QueryMod) (int, int) {
	// group to remove duplicates
	if p.GroupBy == "" {
		*mods = append(*mods, qm.GroupBy("id"))
	} else {
		*mods = append(*mods, qm.GroupBy(p.GroupBy))
	}

	if p.OrderBy == "" {
		*mods = append(*mods, qm.OrderBy("created_at desc"))
	} else {
		*mods = append(*mods, qm.OrderBy(p.OrderBy))
	}

	var limit, offset int
	if p.PageSize == 0 {
		limit = common.APIDefaultPageSize
	} else {
		limit = mathutil.Min(p.PageSize, common.APIMaxPageSize)
	}
	if p.PageNumber > 1 {
		offset = (p.PageNumber - 1) * limit
	}

	*mods = append(*mods, qm.Limit(limit))
	if offset != 0 {
		*mods = append(*mods, qm.Offset(offset))
	}

	return limit, offset
}

func ParseListParams(r *http.Request) (*ListParams, error) {
	params := &ListParams{
		PageNumber: 1,
		PageSize:   50,
	}

	query := r.URL.Query()

	strVal := query.Get("page_no")
	if strVal != "" {
		if val, err := strconv.Atoi(strVal); err != nil {
			return nil, fmt.Errorf("page_no is not an integer: %w", err)
		} else if val < 1 {
			return nil, fmt.Errorf("page_no must be at least 1")
		} else {
			params.PageNumber = val
		}
	}

	strVal = query.Get("page_size")
	if strVal != "" {
		if val, err := strconv.Atoi(strVal); err != nil {
			return nil, fmt.Errorf("page_size is not an integer: %w", err)
		} else if val < 1 {
			return nil, fmt.Errorf("page_size must be at least 1")
		} else {
			params.PageSize = val
		}
	}

	strVal = query.Get("order_by")
	if strVal != "" {
		params.OrderBy = strVal
	}

	return params, nil
}

type ListResponse struct {
	Total int64 `json:"total"`
}

type RoomsResponse struct {
	ListResponse
	Rooms []*models.Room `json:"data"`
}
