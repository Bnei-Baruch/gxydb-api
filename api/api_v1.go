package api

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/edoshor/janus-go"
	"github.com/gorilla/mux"
	pkgerr "github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/gxydb-api/models"
	"github.com/Bnei-Baruch/gxydb-api/pkg/httputil"
	"github.com/Bnei-Baruch/gxydb-api/pkg/middleware"
	"github.com/Bnei-Baruch/gxydb-api/pkg/sqlutil"
)

type V1User struct {
	ID        string `json:"id"`
	Display   string `json:"display"`
	Email     string `json:"email"`
	Group     string `json:"group"`
	IP        string `json:"ip"`
	Janus     string `json:"janus"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	System    string `json:"system"`
	Username  string `json:"username"`
	Room      int    `json:"room"`
	Timestamp int64  `json:"timestamp"`
	Session   int64  `json:"session"`
	Handle    int64  `json:"handle"`
	RFID      int64  `json:"rfid"`
	Camera    bool   `json:"camera"`
	Question  bool   `json:"question"`
	SelfTest  bool   `json:"self_test"`
	SoundTest bool   `json:"sound_test"`
}

type V1RoomInfo struct {
	Room        int    `json:"room"`
	Janus       string `json:"janus"`
	Description string `json:"description"`
}

type V1Room struct {
	V1RoomInfo
	Questions          bool      `json:"questions"`
	NumUsers           int       `json:"num_users"`
	Users              []*V1User `json:"users"`
	firstSessionInRoom time.Time
}

type V1Composite struct {
	VQuad []*V1CompositeRoom `json:"vquad"`
}

type V1CompositeRoom struct {
	V1Room
	Position int `json:"queue"`
}

type V1ProtocolMessageText struct {
	Type   string
	Status bool
	Room   int
	User   V1User
}

func (a *App) V1ListGroups(w http.ResponseWriter, r *http.Request) {
	rooms := a.cache.rooms.Values()
	roomInfos := make([]*V1RoomInfo, len(rooms))
	for i := range rooms {
		room := rooms[i]

		gateway, ok := a.cache.gateways.ByID(room.DefaultGatewayID)
		if !ok {
			log.Ctx(r.Context()).Error().Msgf("gateways cache miss %d [room %d]", room.DefaultGatewayID, room.ID)
			continue
		}

		roomInfos[i] = &V1RoomInfo{
			Room:        room.GatewayUID,
			Janus:       gateway.Name,
			Description: room.Name,
		}
	}

	sort.SliceStable(roomInfos, func(i, j int) bool {
		return roomInfos[i].Description < roomInfos[j].Description
	})

	httputil.RespondWithJSON(w, http.StatusOK, map[string]interface{}{"rooms": roomInfos})
}

func (a *App) V1CreateGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		httputil.NewBadRequestError(err, "malformed id").Abort(w, r)
		return
	}

	var data V1RoomInfo
	if err := httputil.DecodeJSONBody(w, r, &data); err != nil {
		err.Abort(w, r)
		return
	}
	a.requestContext(r).Params = data

	gateway, ok := a.cache.gateways.ByName(data.Janus)
	if !ok {
		httputil.NewBadRequestError(nil, fmt.Sprintf("unknown gateway %s", data.Janus)).Abort(w, r)
		return
	}

	room := models.Room{
		Name:             data.Description,
		DefaultGatewayID: gateway.ID,
		GatewayUID:       id,
	}

	err = sqlutil.InTx(r.Context(), a.DB, func(tx *sql.Tx) error {
		if err := room.Upsert(tx, true, []string{models.RoomColumns.GatewayUID}, boil.Infer(), boil.Infer()); err != nil {
			return pkgerr.WithStack(err)
		}

		a.cache.rooms.Set(&room)

		return nil
	})

	if err != nil {
		httputil.NewInternalError(err).Abort(w, r)
		return
	}

	httputil.RespondSuccess(w)
}

func (a *App) V1ListRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := models.Rooms(
		models.RoomWhere.Disabled.EQ(false),
		models.RoomWhere.RemovedAt.IsNull(),
		qm.Load(models.RoomRels.Sessions, models.SessionWhere.RemovedAt.IsNull()),
		qm.Load(qm.Rels(models.RoomRels.Sessions, models.SessionRels.User)),
	).All(a.DB)

	if err != nil {
		httputil.NewInternalError(pkgerr.WithStack(err)).Abort(w, r)
		return
	}

	respRooms := make([]*V1Room, 0)
	for _, room := range rooms {

		// TODO: maybe move this into DB query above for performance reasons ?
		if len(room.R.Sessions) == 0 {
			continue
		}

		gateway, ok := a.cache.gateways.ByID(room.DefaultGatewayID)
		if !ok {
			log.Ctx(r.Context()).Error().Msgf("gateways cache miss %d [room %d]", room.DefaultGatewayID, room.ID)
			continue
		}

		respRoom := &V1Room{
			V1RoomInfo: V1RoomInfo{
				Room:        room.GatewayUID,
				Janus:       gateway.Name,
				Description: room.Name,
			},
			NumUsers: len(room.R.Sessions),
			Users:    make([]*V1User, len(room.R.Sessions)),
		}

		for i, session := range room.R.Sessions {
			if session.Question {
				respRoom.Questions = true
			}
			if respRoom.firstSessionInRoom.IsZero() || respRoom.firstSessionInRoom.After(session.CreatedAt) {
				respRoom.firstSessionInRoom = session.CreatedAt
			}
			respRoom.Users[i] = a.makeV1User(room, session)
		}

		respRooms = append(respRooms, respRoom)
	}

	// TODO: maybe move to client ?!
	sort.Slice(respRooms, func(i, j int) bool {
		return respRooms[i].firstSessionInRoom.Before(respRooms[j].firstSessionInRoom)
	})

	httputil.RespondWithJSON(w, http.StatusOK, respRooms)
}

func (a *App) V1GetRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		httputil.NewBadRequestError(err, "malformed id").Abort(w, r)
		return
	}

	cachedRoom, ok := a.cache.rooms.ByGatewayUID(id)
	if !ok {
		httputil.NewNotFoundError().Abort(w, r)
		return
	}

	room, err := models.Rooms(
		models.RoomWhere.ID.EQ(cachedRoom.ID),
		models.RoomWhere.Disabled.EQ(false),
		models.RoomWhere.RemovedAt.IsNull(),
		qm.Load(models.RoomRels.DefaultGateway),
		qm.Load(models.RoomRels.Sessions, models.SessionWhere.RemovedAt.IsNull(), qm.OrderBy(models.SessionColumns.CreatedAt)),
		qm.Load(qm.Rels(models.RoomRels.Sessions, models.SessionRels.User)),
		qm.Load(qm.Rels(models.RoomRels.Sessions, models.SessionRels.Gateway)),
	).One(a.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httputil.NewNotFoundError().Abort(w, r)
		} else {
			httputil.NewInternalError(pkgerr.WithStack(err)).Abort(w, r)
		}
		return
	}

	respRoom := &V1Room{
		V1RoomInfo: V1RoomInfo{
			Room:        room.GatewayUID,
			Janus:       room.R.DefaultGateway.Name,
			Description: room.Name,
		},
		NumUsers: len(room.R.Sessions),
		Users:    make([]*V1User, len(room.R.Sessions)),
	}

	for i, session := range room.R.Sessions {
		if session.Question {
			respRoom.Questions = true
		}
		respRoom.Users[i] = a.makeV1User(room, session)
	}

	httputil.RespondWithJSON(w, http.StatusOK, respRoom)
}

func (a *App) V1ListUsers(w http.ResponseWriter, r *http.Request) {
	sessions, err := models.Sessions(
		models.SessionWhere.RemovedAt.IsNull(),
		qm.Load(models.SessionRels.User),
		qm.Load(models.SessionRels.Room),
	).All(a.DB)

	if err != nil {
		httputil.NewInternalError(pkgerr.WithStack(err)).Abort(w, r)
		return
	}

	respSessions := make(map[string]*V1User, len(sessions))
	for i := range sessions {
		session := sessions[i]
		respSessions[session.R.User.AccountsID] = a.makeV1User(session.R.Room, session)
	}

	httputil.RespondWithJSON(w, http.StatusOK, respSessions)
}

func (a *App) V1GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) > 36 {
		httputil.NewBadRequestError(nil, "malformed id").Abort(w, r)
		return
	}

	var userID int64
	err := models.Users(
		qm.Select(models.UserColumns.ID),
		models.UserWhere.AccountsID.EQ(id),
		models.UserWhere.Disabled.EQ(false),
		models.UserWhere.RemovedAt.IsNull(),
	).QueryRow(a.DB).Scan(&userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httputil.NewNotFoundError().Abort(w, r)
		} else {
			httputil.NewInternalError(pkgerr.WithStack(err)).Abort(w, r)
		}
		return
	}

	session, err := models.Sessions(
		models.SessionWhere.UserID.EQ(userID),
		models.SessionWhere.RemovedAt.IsNull(),
		qm.Load(models.SessionRels.User),
		qm.Load(models.SessionRels.Room),
	).One(a.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httputil.NewNotFoundError().Abort(w, r)
		} else {
			httputil.NewInternalError(pkgerr.WithStack(err)).Abort(w, r)
		}
		return
	}

	httputil.RespondWithJSON(w, http.StatusOK, a.makeV1User(session.R.Room, session))
}

func (a *App) V1ListComposites(w http.ResponseWriter, r *http.Request) {
	composites, err := models.Composites(
		qm.Load(models.CompositeRels.CompositesRooms, qm.OrderBy(models.CompositesRoomColumns.Position)),
		qm.Load(qm.Rels(models.CompositeRels.CompositesRooms, models.CompositesRoomRels.Room),
			models.RoomWhere.Disabled.EQ(false),
			models.RoomWhere.RemovedAt.IsNull()),
		qm.Load(qm.Rels(models.CompositeRels.CompositesRooms, models.CompositesRoomRels.Room, models.RoomRels.Sessions),
			models.SessionWhere.RemovedAt.IsNull()),
	).All(a.DB)

	if err != nil {
		httputil.NewInternalError(pkgerr.WithStack(err)).Abort(w, r)
		return
	}

	respComposites := make(map[string]*V1Composite, len(composites))
	for _, composite := range composites {
		respComposites[composite.Name] = a.makeV1Composite(composite)
	}

	httputil.RespondWithJSON(w, http.StatusOK, respComposites)
}

func (a *App) V1GetComposite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) > 16 {
		httputil.NewBadRequestError(nil, "malformed id").Abort(w, r)
		return
	}

	composite, err := models.Composites(
		models.CompositeWhere.Name.EQ(id),
		qm.Load(models.CompositeRels.CompositesRooms),
		qm.Load(qm.Rels(models.CompositeRels.CompositesRooms, models.CompositesRoomRels.Room),
			models.RoomWhere.Disabled.EQ(false),
			models.RoomWhere.RemovedAt.IsNull()),
		qm.Load(qm.Rels(models.CompositeRels.CompositesRooms, models.CompositesRoomRels.Room, models.RoomRels.Sessions),
			models.SessionWhere.RemovedAt.IsNull(), qm.OrderBy(models.SessionColumns.CreatedAt)),
	).One(a.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httputil.NewNotFoundError().Abort(w, r)
		} else {
			httputil.NewInternalError(pkgerr.WithStack(err)).Abort(w, r)
		}
		return
	}

	respComposite := a.makeV1Composite(composite)
	httputil.RespondWithJSON(w, http.StatusOK, respComposite)
}

func (a *App) V1UpdateComposite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) > 16 {
		httputil.NewBadRequestError(nil, "malformed id").Abort(w, r)
		return
	}

	var data V1Composite
	if err := httputil.DecodeJSONBody(w, r, &data); err != nil {
		err.Abort(w, r)
		return
	}
	a.requestContext(r).Params = data

	composite, err := models.Composites(
		models.CompositeWhere.Name.EQ(id),
	).One(a.DB)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httputil.NewNotFoundError().Abort(w, r)
		} else {
			httputil.NewInternalError(pkgerr.WithStack(err)).Abort(w, r)
		}
		return
	}

	err = sqlutil.InTx(r.Context(), a.DB, func(tx *sql.Tx) error {
		cRooms := make(models.CompositesRoomSlice, len(data.VQuad))
		for i, item := range data.VQuad {
			gateway, ok := a.cache.gateways.ByName(item.Janus)
			if !ok {
				return httputil.NewBadRequestError(nil, fmt.Sprintf("unknown gateway %s", item.Janus))
			}
			room, ok := a.cache.rooms.ByGatewayUID(item.Room)
			if !ok {
				return httputil.NewBadRequestError(nil, fmt.Sprintf("unknown room %d", item.Room))
			}

			cRooms[i] = &models.CompositesRoom{
				RoomID:    room.ID,
				GatewayID: gateway.ID,
				Position:  i + 1,
			}
		}

		if _, err := composite.CompositesRooms().DeleteAll(tx); err != nil {
			return pkgerr.WithStack(err)
		}
		if err := composite.AddCompositesRooms(tx, true, cRooms...); err != nil {
			return pkgerr.WithStack(err)
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

	httputil.RespondSuccess(w)
}

func (a *App) V1HandleEvent(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httputil.NewBadRequestError(err, "read request body failed").Abort(w, r)
		return
	}
	r.Body.Close()
	rCtx := a.requestContext(r)
	rCtx.Params = body

	event, err := janus.ParseEvent(body)
	if err != nil {
		httputil.NewBadRequestError(err, "error parsing request body").Abort(w, r)
		return
	}
	rCtx.Params = event

	if err := a.sessionManager.HandleEvent(r.Context(), event); err != nil {
		httputil.NewInternalError(err).Abort(w, r)
		return
	}

	httputil.RespondSuccess(w)
}

func (a *App) V1HandleProtocol(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httputil.NewBadRequestError(err, "read request body failed").Abort(w, r)
		return
	}
	r.Body.Close()
	rCtx := a.requestContext(r)
	rCtx.Params = body

	msg, err := janus.ParseTextroomMessage(body)
	if err != nil {
		httputil.NewBadRequestError(err, "error parsing request body").Abort(w, r)
		return
	}
	rCtx.Params = msg

	if err := a.sessionManager.HandleProtocol(r.Context(), msg); err != nil {
		var pErr *ProtocolError
		if errors.As(err, &pErr) {
			httputil.NewBadRequestError(err, "protocol error").Abort(w, r)
		} else {
			httputil.NewInternalError(err).Abort(w, r)
		}
		return
	}

	httputil.RespondSuccess(w)
}

func (a *App) makeV1User(room *models.Room, session *models.Session) *V1User {
	user := &V1User{
		ID:        session.R.User.AccountsID,
		Display:   session.Display.String,
		Email:     session.R.User.Email.String,
		Group:     room.Name,
		IP:        session.IPAddress.String,
		Name:      "",     // Useless. Shouldn't be used on the client side.
		Role:      "user", // fixed. No more "groups" only "users"
		System:    session.UserAgent.String,
		Username:  session.R.User.Username.String, // Useless. Never seen a value here
		Room:      room.GatewayUID,
		Timestamp: session.CreatedAt.Unix(), // Not sure we really need this
		Session:   session.GatewaySession.Int64,
		Handle:    session.GatewayHandle.Int64,
		RFID:      session.GatewayFeed.Int64,
		Camera:    session.Camera,
		Question:  session.Question,
		SelfTest:  session.SelfTest,  // Not sure we really need this
		SoundTest: session.SoundTest, // Not sure we really need this
	}

	if session.GatewayID.Valid {
		gateway, _ := a.cache.gateways.ByID(session.GatewayID.Int64)
		user.Janus = gateway.Name
	}

	return user
}

func (a *App) makeV1Composite(composite *models.Composite) *V1Composite {
	respComposite := &V1Composite{
		VQuad: make([]*V1CompositeRoom, len(composite.R.CompositesRooms)),
	}

	for j, cRoom := range composite.R.CompositesRooms {
		room := cRoom.R.Room
		gateway, _ := a.cache.gateways.ByID(cRoom.GatewayID)
		respRoom := &V1CompositeRoom{
			V1Room: V1Room{
				V1RoomInfo: V1RoomInfo{
					Room:        room.GatewayUID,
					Janus:       gateway.Name,
					Description: room.Name,
				},
				NumUsers: len(room.R.Sessions),
			},
			Position: cRoom.Position,
		}

		for _, session := range room.R.Sessions {
			if session.Question {
				respRoom.Questions = true
				break
			}
		}

		respComposite.VQuad[j] = respRoom
	}

	return respComposite
}

func (a *App) requestContext(r *http.Request) *middleware.RequestContext {
	rCtx, _ := middleware.ContextFromRequest(r)
	return rCtx
}
