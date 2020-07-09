package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"

	janus_plugins "github.com/edoshor/janus-go/plugins"

	"github.com/Bnei-Baruch/gxydb-api/common"
	"github.com/Bnei-Baruch/gxydb-api/domain"
	"github.com/Bnei-Baruch/gxydb-api/models"
	"github.com/Bnei-Baruch/gxydb-api/pkg/stringutil"
)

func (s *ApiTestSuite) TestAdmin_GatewaysHandleInfoForbidden() {
	req, _ := http.NewRequest("GET", "/admin/gateways/1/sessions/1/handles/1/info", nil)
	resp := s.request(req)
	s.Require().Equal(http.StatusUnauthorized, resp.Code)

	req, _ = http.NewRequest("GET", "/admin/gateways/1/sessions/1/handles/1/info", nil)
	s.apiAuth(req)
	resp = s.request(req)
	s.Require().Equal(http.StatusForbidden, resp.Code)
}

func (s *ApiTestSuite) TestAdmin_GatewaysHandleInfoNotFound() {
	req, _ := http.NewRequest("GET", "/admin/gateways/1/sessions/1/handles/1/info", nil)
	s.apiAuthP(req, []string{common.RoleAdmin})
	resp := s.request(req)
	s.Require().Equal(http.StatusNotFound, resp.Code)

	gateway := s.createGatewayP(common.GatewayTypeRooms, s.GatewayManager.Config.AdminURL, s.GatewayManager.Config.AdminSecret)
	s.Require().NoError(s.app.cache.ReloadAll(s.DB))

	req, _ = http.NewRequest("GET", fmt.Sprintf("/admin/gateways/%s/sessions/1/handles/1/info", gateway.Name), nil)
	s.apiAuthP(req, []string{common.RoleAdmin})
	resp = s.request(req)
	s.Require().Equal(http.StatusNotFound, resp.Code)

	session, err := s.NewGatewaySession()
	s.Require().NoError(err, "NewGatewaySession")
	defer session.Destroy()

	req, _ = http.NewRequest("GET", fmt.Sprintf("/admin/gateways/%s/sessions/%d/handles/1/info", gateway.Name, session.Id), nil)
	s.apiAuthP(req, []string{common.RoleAdmin})
	resp = s.request(req)
	s.Require().Equal(http.StatusNotFound, resp.Code)
}

func (s *ApiTestSuite) TestAdmin_GatewaysHandleInfo() {
	gateway := s.createGatewayP(common.GatewayTypeRooms, s.GatewayManager.Config.AdminURL, s.GatewayManager.Config.AdminSecret)
	s.Require().NoError(s.app.cache.ReloadAll(s.DB))

	session, err := s.NewGatewaySession()
	s.Require().NoError(err, "NewGatewaySession")
	defer session.Destroy()

	handle, err := session.Attach("janus.plugin.videoroom")
	s.Require().NoError(err, "session.Attach")
	defer handle.Detach()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/admin/gateways/%s/sessions/%d/handles/%d/info", gateway.Name, session.Id, handle.Id), nil)
	s.apiAuthP(req, []string{common.RoleAdmin})
	body := s.request200json(req)
	s.EqualValues(session.Id, body["session_id"], "session_id")
	s.EqualValues(handle.Id, body["handle_id"], "handle_id")
	s.NotNil(body["info"], "info")
}

func (s *ApiTestSuite) TestAdmin_ListRoomsForbidden() {
	req, _ := http.NewRequest("GET", "/admin/rooms", nil)
	resp := s.request(req)
	s.Require().Equal(http.StatusUnauthorized, resp.Code)

	req, _ = http.NewRequest("GET", "/admin/rooms", nil)
	s.apiAuth(req)
	resp = s.request(req)
	s.Require().Equal(http.StatusForbidden, resp.Code)
}

func (s *ApiTestSuite) TestAdmin_ListRoomsBadRequest() {
	args := [...]string{
		"page_no=0",
		"page_no=-1",
		"page_no=abc",
		"page_size=0",
		"page_size=-1",
		"page_size=abc",
	}
	for i, query := range args {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/admin/rooms?%s", query), nil)
		s.apiAuthP(req, []string{common.RoleRoot})
		resp := s.request(req)
		s.Require().Equal(http.StatusBadRequest, resp.Code, i)
	}
}

func (s *ApiTestSuite) TestAdmin_ListRooms() {
	req, _ := http.NewRequest("GET", "/admin/rooms", nil)
	s.apiAuthP(req, []string{common.RoleRoot})
	body := s.request200json(req)
	s.Equal(0, int(body["total"].(float64)), "total")
	s.Equal(0, len(body["data"].([]interface{})), "len(data)")

	gateway := s.createGateway()
	rooms := make([]*models.Room, 10)
	for i := range rooms {
		rooms[i] = s.createRoom(gateway)
	}

	body = s.request200json(req)
	s.Equal(10, int(body["total"].(float64)), "total")
	s.Equal(10, len(body["data"].([]interface{})), "len(data)")

	for i, room := range rooms {
		req, _ = http.NewRequest("GET", fmt.Sprintf("/admin/rooms?page_no=%d&page_size=1&order_by=id", i+1), nil)
		s.apiAuthP(req, []string{common.RoleRoot})
		body = s.request200json(req)

		s.Equal(10, int(body["total"].(float64)), "total")
		data := body["data"].([]interface{})
		s.Equal(1, len(data), "len(data)")
		roomData := data[0].(map[string]interface{})
		s.Equal(roomData["name"], room.Name, "name")
		s.EqualValues(roomData["default_gateway_id"], room.DefaultGatewayID, "default_gateway_id")
		s.EqualValues(roomData["gateway_uid"], room.GatewayUID, "gateway_uid")
	}
}

func (s *ApiTestSuite) TestAdmin_GetRoomForbidden() {
	req, _ := http.NewRequest("GET", "/admin/rooms/1", nil)
	resp := s.request(req)
	s.Require().Equal(http.StatusUnauthorized, resp.Code)

	req, _ = http.NewRequest("GET", "/admin/rooms/1", nil)
	s.apiAuth(req)
	resp = s.request(req)
	s.Require().Equal(http.StatusForbidden, resp.Code)
}

func (s *ApiTestSuite) TestAdmin_GetRoomNotFound() {
	req, _ := http.NewRequest("GET", "/admin/rooms/abc", nil)
	s.apiAuthP(req, []string{common.RoleRoot})
	resp := s.request(req)
	s.Require().Equal(http.StatusNotFound, resp.Code)

	req, _ = http.NewRequest("GET", "/admin/rooms/1", nil)
	s.apiAuthP(req, []string{common.RoleRoot})
	resp = s.request(req)
	s.Require().Equal(http.StatusNotFound, resp.Code)
}

func (s *ApiTestSuite) TestAdmin_GetRoom() {
	gateway := s.createGateway()
	room := s.createRoom(gateway)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/admin/rooms/%d", room.ID), nil)
	s.apiAuthP(req, []string{common.RoleRoot})
	body := s.request200json(req)
	s.Equal(body["name"], room.Name, "name")
	s.EqualValues(body["default_gateway_id"], room.DefaultGatewayID, "default_gateway_id")
	s.EqualValues(body["gateway_uid"], room.GatewayUID, "gateway_uid")
}

func (s *ApiTestSuite) TestAdmin_CreateRoomForbidden() {
	req, _ := http.NewRequest("POST", "/admin/rooms", nil)
	resp := s.request(req)
	s.Require().Equal(http.StatusUnauthorized, resp.Code)

	req, _ = http.NewRequest("POST", "/admin/rooms", nil)
	s.apiAuth(req)
	resp = s.request(req)
	s.Require().Equal(http.StatusForbidden, resp.Code)
}

func (s *ApiTestSuite) TestAdmin_CreateRoomBadRequest() {
	req, _ := http.NewRequest("POST", "/admin/rooms", bytes.NewBuffer([]byte("{\"bad\":\"json")))
	s.apiAuthP(req, []string{common.RoleRoot})
	resp := s.request(req)
	s.Require().Equal(http.StatusBadRequest, resp.Code)

	// non existing gateway
	body := models.Room{
		Name:       fmt.Sprintf("room_%s", stringutil.GenerateName(10)),
		GatewayUID: rand.Intn(math.MaxInt32),
	}
	b, _ := json.Marshal(body)
	req, _ = http.NewRequest("POST", "/admin/rooms", bytes.NewBuffer(b))
	s.apiAuthP(req, []string{common.RoleRoot})
	resp = s.request(req)
	s.Require().Equal(http.StatusBadRequest, resp.Code)

	// invalid gateway uid
	gateway := s.createGatewayP(common.GatewayTypeRooms, s.GatewayManager.Config.AdminURL, s.GatewayManager.Config.AdminSecret)
	s.Require().NoError(s.app.cache.ReloadAll(s.DB))

	body.DefaultGatewayID = gateway.ID
	body.GatewayUID = -8
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest("POST", "/admin/rooms", bytes.NewBuffer(b))
	s.apiAuthP(req, []string{common.RoleRoot})
	resp = s.request(req)
	s.Require().Equal(http.StatusBadRequest, resp.Code)

	// existing gateway_uid
	room := s.createRoom(gateway)
	s.Require().NoError(s.app.cache.ReloadAll(s.DB))
	body.GatewayUID = room.GatewayUID
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest("POST", "/admin/rooms", bytes.NewBuffer(b))
	s.apiAuthP(req, []string{common.RoleRoot})
	resp = s.request(req)
	s.Require().Equal(http.StatusBadRequest, resp.Code)

	// existing name
	body.Name = room.Name
	body.GatewayUID = room.GatewayUID + 1
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest("POST", "/admin/rooms", bytes.NewBuffer(b))
	s.apiAuthP(req, []string{common.RoleRoot})
	resp = s.request(req)
	s.Require().Equal(http.StatusBadRequest, resp.Code)

	// invalid name
	body.GatewayUID = room.GatewayUID
	for _, name := range []string{"", "אסור עברית", "123456789012345678901234567890123456789012345678901234567890123456789012345"} {
		body.Name = name
		b, _ = json.Marshal(body)
		req, _ = http.NewRequest("POST", "/admin/rooms", bytes.NewBuffer(b))
		s.apiAuthP(req, []string{common.RoleRoot})
		resp = s.request(req)
		s.Require().Equal(http.StatusBadRequest, resp.Code)
	}
}

func (s *ApiTestSuite) TestAdmin_CreateRoom() {
	gateway := s.createGatewayP(common.GatewayTypeRooms, s.GatewayManager.Config.AdminURL, s.GatewayManager.Config.AdminSecret)
	s.Require().NoError(s.app.cache.ReloadAll(s.DB))

	payload := models.Room{
		Name:             fmt.Sprintf("room_%s", stringutil.GenerateName(10)),
		GatewayUID:       rand.Intn(math.MaxInt32),
		DefaultGatewayID: gateway.ID,
	}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/admin/rooms", bytes.NewBuffer(b))
	s.apiAuthP(req, []string{common.RoleRoot})
	body := s.request201json(req)
	s.NotZero(body["id"], "id")
	s.Equal(payload.Name, body["name"], "name")
	s.EqualValues(payload.GatewayUID, body["gateway_uid"], "gateway_uid")
	s.EqualValues(payload.DefaultGatewayID, body["default_gateway_id"], "default_gateway_id")
	s.False(body["disabled"].(bool), "disabled")

	// verify room is created on gateway
	gRoom := s.findRoomInGateway(gateway, int(body["gateway_uid"].(float64)))
	s.Require().NotNil(gRoom, "gateway room")
	s.Equal(gRoom.Description, payload.Name, "gateway room description")
}

func (s *ApiTestSuite) TestAdmin_UpdateRoomForbidden() {
	req, _ := http.NewRequest("PUT", "/admin/rooms/1", nil)
	resp := s.request(req)
	s.Require().Equal(http.StatusUnauthorized, resp.Code)

	req, _ = http.NewRequest("PUT", "/admin/rooms/1", nil)
	s.apiAuth(req)
	resp = s.request(req)
	s.Require().Equal(http.StatusForbidden, resp.Code)
}

func (s *ApiTestSuite) TestAdmin_UpdateRoomNotFound() {
	req, _ := http.NewRequest("PUT", "/admin/rooms/abc", nil)
	s.apiAuthP(req, []string{common.RoleRoot})
	resp := s.request(req)
	s.Require().Equal(http.StatusNotFound, resp.Code)

	req, _ = http.NewRequest("PUT", "/admin/rooms/1", nil)
	s.apiAuthP(req, []string{common.RoleRoot})
	resp = s.request(req)
	s.Require().Equal(http.StatusNotFound, resp.Code)
}

func (s *ApiTestSuite) TestAdmin_UpdateRoomBadRequest() {
	gateway := s.createGatewayP(common.GatewayTypeRooms, s.GatewayManager.Config.AdminURL, s.GatewayManager.Config.AdminSecret)
	room := s.createRoom(gateway)
	s.Require().NoError(s.app.cache.ReloadAll(s.DB))

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/admin/rooms/%d", room.ID), bytes.NewBuffer([]byte("{\"bad\":\"json")))
	s.apiAuthP(req, []string{common.RoleRoot})
	resp := s.request(req)
	s.Require().Equal(http.StatusBadRequest, resp.Code)

	// non existing gateway
	body := models.Room{
		Name:       fmt.Sprintf("room_%s", stringutil.GenerateName(10)),
		GatewayUID: rand.Intn(math.MaxInt32),
	}
	b, _ := json.Marshal(body)
	req, _ = http.NewRequest("PUT", fmt.Sprintf("/admin/rooms/%d", room.ID), bytes.NewBuffer(b))
	s.apiAuthP(req, []string{common.RoleRoot})
	resp = s.request(req)
	s.Require().Equal(http.StatusBadRequest, resp.Code)

	// invalid gateway uid
	body.DefaultGatewayID = gateway.ID
	body.GatewayUID = -8
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest("PUT", fmt.Sprintf("/admin/rooms/%d", room.ID), bytes.NewBuffer(b))
	s.apiAuthP(req, []string{common.RoleRoot})
	resp = s.request(req)
	s.Require().Equal(http.StatusBadRequest, resp.Code)

	// existing gateway_uid
	room2 := s.createRoom(gateway)
	s.Require().NoError(s.app.cache.ReloadAll(s.DB))
	body.GatewayUID = room2.GatewayUID
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest("PUT", fmt.Sprintf("/admin/rooms/%d", room.ID), bytes.NewBuffer(b))
	s.apiAuthP(req, []string{common.RoleRoot})
	resp = s.request(req)
	s.Require().Equal(http.StatusBadRequest, resp.Code)

	// existing name
	body.GatewayUID = room.GatewayUID
	body.Name = room2.Name
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest("PUT", fmt.Sprintf("/admin/rooms/%d", room.ID), bytes.NewBuffer(b))
	s.apiAuthP(req, []string{common.RoleRoot})
	resp = s.request(req)
	s.Require().Equal(http.StatusBadRequest, resp.Code)

	// invalid name
	for _, name := range []string{"", "אסור עברית", "123456789012345678901234567890123456789012345678901234567890123456789012345"} {
		body.Name = name
		b, _ = json.Marshal(body)
		req, _ = http.NewRequest("POST", "/admin/rooms", bytes.NewBuffer(b))
		s.apiAuthP(req, []string{common.RoleRoot})
		resp = s.request(req)
		s.Require().Equal(http.StatusBadRequest, resp.Code)
	}
}

func (s *ApiTestSuite) TestAdmin_UpdateRoom() {
	gateway := s.createGatewayP(common.GatewayTypeRooms, s.GatewayManager.Config.AdminURL, s.GatewayManager.Config.AdminSecret)
	s.Require().NoError(s.app.cache.ReloadAll(s.DB))

	payload := models.Room{
		Name:             fmt.Sprintf("room_%s", stringutil.GenerateName(10)),
		GatewayUID:       rand.Intn(math.MaxInt16),
		DefaultGatewayID: gateway.ID,
	}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/admin/rooms", bytes.NewBuffer(b))
	s.apiAuthP(req, []string{common.RoleRoot})
	body := s.request201json(req)

	gateway2 := s.createGatewayP(common.GatewayTypeRooms, s.GatewayManager.Config.AdminURL, s.GatewayManager.Config.AdminSecret)
	s.Require().NoError(s.app.cache.ReloadAll(s.DB))

	payload.Name = fmt.Sprintf("%s_edit", payload.Name)
	payload.DefaultGatewayID = gateway2.ID
	payload.Disabled = true
	b, _ = json.Marshal(payload)
	req, _ = http.NewRequest("PUT", fmt.Sprintf("/admin/rooms/%d", int64(body["id"].(float64))), bytes.NewBuffer(b))
	s.apiAuthP(req, []string{common.RoleRoot})
	body = s.request200json(req)
	s.Equal(payload.Name, body["name"], "name")
	s.EqualValues(payload.DefaultGatewayID, body["default_gateway_id"], "default_gateway_id")
	s.True(body["disabled"].(bool), "disabled")
	s.Greater(body["updated_at"], body["created_at"], "updated_at > created_at")

	// verify room is updated on gateway
	gRoom := s.findRoomInGateway(gateway, int(body["gateway_uid"].(float64)))
	s.Require().NotNil(gRoom, "gateway room")
	s.Equal(gRoom.Description, payload.Name, "gateway room description")
}

func (s *ApiTestSuite) TestAdmin_DeleteRoomForbidden() {
	req, _ := http.NewRequest("DELETE", "/admin/rooms/1", nil)
	resp := s.request(req)
	s.Require().Equal(http.StatusUnauthorized, resp.Code)

	req, _ = http.NewRequest("DELETE", "/admin/rooms/1", nil)
	s.apiAuth(req)
	resp = s.request(req)
	s.Require().Equal(http.StatusForbidden, resp.Code)
}

func (s *ApiTestSuite) TestAdmin_DeleteRoomNotFound() {
	req, _ := http.NewRequest("DELETE", "/admin/rooms/abc", nil)
	s.apiAuthP(req, []string{common.RoleRoot})
	resp := s.request(req)
	s.Require().Equal(http.StatusNotFound, resp.Code)

	req, _ = http.NewRequest("DELETE", "/admin/rooms/1", nil)
	s.apiAuthP(req, []string{common.RoleRoot})
	resp = s.request(req)
	s.Require().Equal(http.StatusNotFound, resp.Code)
}

func (s *ApiTestSuite) TestAdmin_DeleteRoom() {
	gateway := s.createGatewayP(common.GatewayTypeRooms, s.GatewayManager.Config.AdminURL, s.GatewayManager.Config.AdminSecret)
	s.Require().NoError(s.app.cache.ReloadAll(s.DB))

	payload := models.Room{
		Name:             fmt.Sprintf("room_%s", stringutil.GenerateName(10)),
		GatewayUID:       rand.Intn(math.MaxInt16),
		DefaultGatewayID: gateway.ID,
	}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/admin/rooms", bytes.NewBuffer(b))
	s.apiAuthP(req, []string{common.RoleRoot})
	body := s.request201json(req)

	id := int64(body["id"].(float64))
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/admin/rooms/%d", id), nil)
	s.apiAuthP(req, []string{common.RoleRoot})
	s.request200json(req)

	// verify room removed_at is set in DB
	room, err := models.FindRoom(s.DB, id)
	s.Require().NoError(err, "models.FindRoom")
	s.True(room.RemovedAt.Valid, "remove_at")

	// verify room does not exists on gateway
	s.Nil(s.findRoomInGateway(gateway, int(body["gateway_uid"].(float64))))
}

func (s *ApiTestSuite) findRoomInGateway(gateway *models.Gateway, id int) *janus_plugins.VideoroomRoomFromListResponse {
	api, err := domain.GatewayAdminAPIRegistry.For(gateway)
	s.Require().NoError(err, "Admin API for gateway")

	request := janus_plugins.MakeVideoroomRequestFactory(common.Config.GatewayVideoroomAdminKey).ListRequest()
	resp, err := api.MessagePlugin(request)
	s.Require().NoError(err, "api.MessagePlugin")

	tResp, _ := resp.(*janus_plugins.VideoroomListResponse)
	for _, x := range tResp.Rooms {
		if x.Room == id {
			return x
		}
	}

	return nil
}
