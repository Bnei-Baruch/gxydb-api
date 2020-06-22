package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	janus_admin "github.com/edoshor/janus-go/admin"
	"github.com/gorilla/mux"
	pkgerr "github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/gxydb-api/common"
	"github.com/Bnei-Baruch/gxydb-api/domain"
	"github.com/Bnei-Baruch/gxydb-api/middleware"
	"github.com/Bnei-Baruch/gxydb-api/models"
	"github.com/Bnei-Baruch/gxydb-api/pkg/httputil"
	"github.com/Bnei-Baruch/gxydb-api/pkg/mathutil"
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

type ListParams struct {
	PageNumber int    `json:"page_no" form:"page_no" binding:"omitempty,min=1"`
	PageSize   int    `json:"page_size" form:"page_size" binding:"omitempty,min=1"`
	OrderBy    string `json:"order_by" form:"order_by" binding:"omitempty"`
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
