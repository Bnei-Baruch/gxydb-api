package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Bnei-Baruch/gxydb-api/common"
	"github.com/Bnei-Baruch/gxydb-api/middleware"
	"github.com/Bnei-Baruch/gxydb-api/pkg/httputil"
)

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

		token, _ := a.cache.gatewayTokens.ByID(gateway.ID)
		respGateway := &V2Gateway{
			Name:  gateway.Name,
			URL:   gateway.URL,
			Type:  gateway.Type,
			Token: token,
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

	data := make(map[int]*V2RoomStatistics, len(stats))
	for _, roomStats := range stats {
		data[roomStats.R.Room.GatewayUID] = &V2RoomStatistics{OnAir: roomStats.OnAir}
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
