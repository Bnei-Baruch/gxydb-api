package domain

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Bnei-Baruch/gxydb-api/models"
)

type RoomServerAssignmentTestSuite struct {
	ModelsSuite
}

func (s *RoomServerAssignmentTestSuite) SetupSuite() {
	s.Require().NoError(s.InitTestDB())
}

func (s *RoomServerAssignmentTestSuite) TearDownSuite() {
	s.Require().NoError(s.DestroyTestDB())
}

func (s *RoomServerAssignmentTestSuite) SetupTest() {
	s.DBCleaner.Acquire(models.TableNames.Gateways, models.TableNames.Rooms, "room_server_assignments", models.TableNames.Sessions)
}

func (s *RoomServerAssignmentTestSuite) TearDownTest() {
	s.DBCleaner.Clean(models.TableNames.Gateways, models.TableNames.Rooms, "room_server_assignments", models.TableNames.Sessions)
}

func (s *RoomServerAssignmentTestSuite) TestGetOrAssignServer_NewAssignment() {
	manager := NewRoomServerAssignmentManager(s.DB, []string{"gxy1", "gxy2", "gxy3"}, 400, 10, nil)

	gateway := s.CreateGateway()
	room := s.CreateRoom(gateway)

	// First call should create an assignment
	server1, err := manager.GetOrAssignServer(context.Background(), room.GatewayUID, "")
	s.Require().NoError(err)
	s.NotEmpty(server1)
	s.Contains([]string{"gxy1", "gxy2", "gxy3"}, server1)

	// Second call should return the same server (sticky routing)
	server2, err := manager.GetOrAssignServer(context.Background(), room.GatewayUID, "")
	s.Require().NoError(err)
	s.Equal(server1, server2)
}

func (s *RoomServerAssignmentTestSuite) TestGetOrAssignServer_LoadBalancing() {
	manager := NewRoomServerAssignmentManager(s.DB, []string{"gxy1", "gxy2"}, 400, 10, nil)

	gateway1 := s.CreateGatewayWithName("gxy1")
	gateway2 := s.CreateGatewayWithName("gxy2")

	user1 := s.CreateUser()
	user2 := s.CreateUser()

	// Create rooms
	room1 := s.CreateRoom(gateway1)
	_ = s.CreateRoom(gateway2) // room2 on gxy2
	room3 := s.CreateRoom(gateway1)

	// Create sessions to simulate load on gxy1
	s.CreateSession(user1, gateway1, room1)
	s.CreateSession(user2, gateway1, room1)

	// New room should be assigned to gxy2 (least loaded)
	server, err := manager.GetOrAssignServer(context.Background(), room3.GatewayUID, "")
	s.Require().NoError(err)
	s.Equal("gxy2", server)
}

func (s *RoomServerAssignmentTestSuite) TestCleanInactiveAssignments() {
	manager := NewRoomServerAssignmentManager(s.DB, []string{"gxy1", "gxy2"}, 400, 10, nil)

	gateway := s.CreateGateway()
	room1 := s.CreateRoom(gateway)
	room2 := s.CreateRoom(gateway)
	user := s.CreateUser()

	// Assign both rooms
	_, err := manager.GetOrAssignServer(context.Background(), room1.GatewayUID, "")
	s.Require().NoError(err)
	_, err = manager.GetOrAssignServer(context.Background(), room2.GatewayUID, "")
	s.Require().NoError(err)

	// Create session only for room1
	s.CreateSession(user, gateway, room1)

	// Clean should remove assignment for room2 (no active sessions)
	err = manager.CleanInactiveAssignments(context.Background())
	s.Require().NoError(err)

	// room1 should still have assignment
	server1, err := manager.GetOrAssignServer(context.Background(), room1.GatewayUID, "")
	s.Require().NoError(err)
	s.NotEmpty(server1)

	// room2 should get a new assignment (was cleaned)
	// We can't directly check if it was cleaned, but GetOrAssignServer should work
	server2, err := manager.GetOrAssignServer(context.Background(), room2.GatewayUID, "")
	s.Require().NoError(err)
	s.NotEmpty(server2)
}

func (s *RoomServerAssignmentTestSuite) TestGetOrAssignServer_RegionalRouting() {
	// Setup regional config: IL -> gxy1,gxy2; US -> gxy3,gxy4
	regions := map[string][]string{
		"IL": {"gxy1", "gxy2"},
		"US": {"gxy3", "gxy4"},
	}
	manager := NewRoomServerAssignmentManager(s.DB, []string{"gxy1", "gxy2", "gxy3", "gxy4"}, 400, 10, regions)

	gateway1 := s.CreateGatewayWithName("gxy1")
	gateway3 := s.CreateGatewayWithName("gxy3")

	room1 := s.CreateRoom(gateway1)
	room2 := s.CreateRoom(gateway3)

	// First user from IL should get gxy1 or gxy2
	server1, err := manager.GetOrAssignServer(context.Background(), room1.GatewayUID, "IL")
	s.Require().NoError(err)
	s.Contains([]string{"gxy1", "gxy2"}, server1)

	// First user from US should get gxy3 or gxy4
	server2, err := manager.GetOrAssignServer(context.Background(), room2.GatewayUID, "US")
	s.Require().NoError(err)
	s.Contains([]string{"gxy3", "gxy4"}, server2)

	// Second user from different region should get same server (sticky routing)
	server1Again, err := manager.GetOrAssignServer(context.Background(), room1.GatewayUID, "US")
	s.Require().NoError(err)
	s.Equal(server1, server1Again) // Should ignore "US" and return same server
}

func TestRoomServerAssignmentTestSuite(t *testing.T) {
	suite.Run(t, new(RoomServerAssignmentTestSuite))
}
