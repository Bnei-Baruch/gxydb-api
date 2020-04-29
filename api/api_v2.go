package api

import (
	"net/http"

	pkgerr "github.com/pkg/errors"

	"github.com/Bnei-Baruch/gxydb-api/models"
	"github.com/Bnei-Baruch/gxydb-api/pkg/config"
	"github.com/Bnei-Baruch/gxydb-api/pkg/httputil"
)

func (a *App) V2GetConfig(w http.ResponseWriter, r *http.Request) {
	gateways := a.cache.gateways.Values()
	cfg := V2Config{
		Gateways:   make(map[string]map[string]*V2Gateway, len(gateways)),
		IceServers: config.Config.IceServers,
	}

	// TODO: implement gateway types
	// TODO: implement gateway tokens
	cfg.Gateways["rooms"] = make(map[string]*V2Gateway)
	for _, gateway := range gateways {
		cfg.Gateways["rooms"][gateway.Name] = &V2Gateway{
			Name:  gateway.Name,
			URL:   gateway.URL,
			Type:  "rooms",
			Token: "secret",
		}
	}

	httputil.RespondWithJSON(w, http.StatusOK, cfg)
}

func (a *App) ListRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := models.Rooms(
		models.RoomWhere.Disabled.EQ(false),
		models.RoomWhere.RemovedAt.IsNull(),
	).All(a.DB)

	if err != nil {
		httputil.NewInternalError(pkgerr.WithStack(err)).Abort(w, r)
		return
	}

	httputil.RespondWithJSON(w, http.StatusOK, rooms)
}
