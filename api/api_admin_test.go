package api

import (
	"fmt"
	"net/http"

	"github.com/Bnei-Baruch/gxydb-api/common"
	"github.com/Bnei-Baruch/gxydb-api/models"
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
}
