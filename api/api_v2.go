package api

import (
	"net/http"

	pkgerr "github.com/pkg/errors"

	"github.com/Bnei-Baruch/gxydb-api/models"
	"github.com/Bnei-Baruch/gxydb-api/pkg/httputil"
)

type V2Gateway struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Type  string `json:"type"`
	Token string `json:"token"`
}

type V2Config struct {
	Gateways   []*V2Gateway `json:"gateways"`
	IceServers []string     `json:"ice_servers"` // TODO: implement
}

func (a *App) V2GetConfig(w http.ResponseWriter, r *http.Request) {
	gateways := a.cache.gateways.Values()
	config := V2Config{
		Gateways:   make([]*V2Gateway, len(gateways)),
		IceServers: []string{},
	}

	for i := range gateways {
		config.Gateways[i] = &V2Gateway{
			Name:  gateways[i].Name,
			URL:   gateways[i].URL,
			Type:  "rooms",  // TODO: implement
			Token: "secret", // TODO: implement
		}
	}

	httputil.RespondWithJSON(w, http.StatusOK, config)
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
