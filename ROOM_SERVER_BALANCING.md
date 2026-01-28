# Room Server Load Balancing

## Overview

Dynamic load balancing implementation for distributing rooms across available Janus servers (gxy1-12).

## Core Principles

1. **Sticky Routing**: The first user entering a room "reserves" a server. All subsequent users will get the same server.
2. **Load Balancing**: When selecting a server for a new room, the current load (number of active sessions) is taken into account.
3. **Auto-cleanup**: Assignments are automatically removed when there are no active users left in the room.

## API

### POST /v2/room_server

Get the server for connecting to a room.

**Request:**
```json
{
  "room": 2171,
  "geo": {
    "country_code": "IL"
  }
}
```

Fields:
- `room` (required): Room number
- `geo.country_code` (optional): User's country code for regional routing (e.g., "IL", "US", "RU")

**Response:**
```json
{
  "janus": "gxy5"
}
```

**Regional Routing:**
- If `geo.country_code` is provided and matches a configured region, server will be selected from regional pool
- Only applies to **first user** in room (assignment creation)
- All subsequent users get the same server regardless of their country (sticky routing)

**Response codes:**
- `200 OK` - server successfully retrieved
- `404 Not Found` - room not found
- `500 Internal Server Error` - server error

## Configuration

### Environment Variables

```bash
# Enable load balancing mode (default: false)
# false = legacy mode (use room's default gateway)
# true = scale mode (dynamic load balancing with regional routing)
SCALE=true

# Available Janus servers (default: gxy1-gxy12) - only for SCALE=true
AVAILABLE_JANUS_SERVERS=gxy1,gxy2,gxy3,gxy4,gxy5,gxy6,gxy7,gxy8,gxy9,gxy10,gxy11,gxy12

# Maximum users per server (default: 400) - only for SCALE=true
MAX_SERVER_CAPACITY=400

# Average users per room (default: 10) - only for SCALE=true
AVG_ROOM_OCCUPANCY=10

# Regional server mapping (optional) - only for SCALE=true
# Format: "COUNTRY_CODE:server1,server2;COUNTRY_CODE2:server3,server4"
SERVER_REGIONS=IL:gxy1,gxy2,gxy3;US:gxy4,gxy5,gxy6;RU:gxy7,gxy8,gxy9
```

**Defaults:**
- Scale mode: `false` (legacy mode)
- Available servers: gxy1-12
- Max server capacity: 400 users
- Average room occupancy: 10 users
- Regional routing: disabled (no regions configured)

## Operating Modes

### Legacy Mode (SCALE=false, default)

Uses room's pre-configured default gateway from database.

**Behavior:**
- Returns `room.default_gateway` name via cache lookup
- No load balancing
- Ignores `geo.country_code` parameter
- Same server for room every time
- No database writes to `room_server_assignments`
- **Strict validation**: Returns error if gateway not found (ensures proper room configuration)

**Use case:** Existing clients that need stable behavior during migration.

**Error Handling:**
If a room's `default_gateway_id` is not found in the gateway cache, the endpoint will:
1. Log an error with room details
2. Return HTTP 500 error to the client
3. Force administrators to fix room configuration

This strict behavior ensures all rooms are properly configured before migrating to scale mode.

### Scale Mode (SCALE=true)

Dynamic load balancing with optional regional routing.

**Behavior:**
- First user: assigns least loaded server (with regional preference if configured)
- Subsequent users: return same server (sticky routing)
- Uses `room_server_assignments` table
- Respects `geo.country_code` for initial assignment
- Auto-cleanup on session end

**Use case:** New deployment with load balancing and regional optimization.

## Migration Strategy

```bash
# Step 1: Start with legacy mode (existing behavior)
SCALE=false

# Step 2: Test scale mode on staging
SCALE=true

# Step 3: Gradually enable on production
# - Deploy with SCALE=false
# - Monitor and test
# - Switch to SCALE=true when ready
# - Can rollback to SCALE=false instantly if issues arise
```

## Database

### Migration

```bash
# Apply migration
migrate -path migrations -database "postgres://user:pass@localhost/galaxy?sslmode=disable" up
```

### Table room_server_assignments

```sql
CREATE TABLE room_server_assignments (
    room_id       BIGINT PRIMARY KEY REFERENCES rooms,
    gateway_name  VARCHAR(50) NOT NULL,
    region        VARCHAR(50),           -- for future functionality
    assigned_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    last_used_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## Algorithm

1. **Server request:**
   - Client sends request with room number
   - Room existence is verified
   - Existing assignment is checked

2. **Existing assignment:**
   - Server name is returned
   - `last_used_at` is updated

3. **New assignment:**
   - Load is calculated for each server (number of active sessions)
   - Server with minimum load is selected
   - Record is created in `room_server_assignments`
   - Selected server name is returned

4. **Auto-cleanup:**
   - Periodically (on session cleanup timer) assignments are removed for rooms without active users
   - On each session update, `last_used_at` is updated for the room assignment

## Integration with Existing Functionality

### MQTT Events

When receiving `leaving`, `kicked`, `unpublished` events via MQTT:
- Session is marked as removed
- Periodic cleaner checks rooms without active sessions
- Removes inactive assignments

### Keepalive

Uses existing `PeriodicSessionCleaner` mechanism:
- Period: `CLEAN_SESSIONS_INTERVAL` (default 1 minute)
- Dead session criteria: `DEAD_SESSION_PERIOD` (default 90 seconds)

## Limitations

- **Server priorities**: Not implemented in current version

## Future Improvements

1. **Server priorities**:
   - Ability to set preferred servers
   - Fallback to other servers on overload

3. **Monitoring**:
   - Server load metrics
   - Room distribution statistics

## Testing

```bash
# Run tests
go test ./domain -run TestRoomServerAssignment
```

## Usage Examples

### Legacy Mode (SCALE=false)

```bash
# Configuration: SCALE=false (or not set)

# Request server for room 2171
curl -X POST http://localhost:8081/v2/room_server \
  -H "Content-Type: application/json" \
  -d '{"room": 2171}'

# Response: {"janus": "gxy5"}  (room's default_gateway from database)

# All users always get the same server
curl -X POST http://localhost:8081/v2/room_server \
  -H "Content-Type: application/json" \
  -d '{"room": 2171, "geo": {"country_code": "US"}}'

# Response: {"janus": "gxy5"}  (same - geo.country_code ignored)
```

### Scale Mode - Basic (SCALE=true, no regions)

```bash
# Configuration: SCALE=true

# 1. First user requests server for room 2171
curl -X POST http://localhost:8081/v2/room_server \
  -H "Content-Type: application/json" \
  -d '{"room": 2171}'

# Response: {"janus": "gxy5"}  (least loaded server selected)

# 2. Second user requests server for same room
curl -X POST http://localhost:8081/v2/room_server \
  -H "Content-Type: application/json" \
  -d '{"room": 2171}'

# Response: {"janus": "gxy5"}  (same server - sticky routing!)
```

### Scale Mode - With Regional Routing (SCALE=true)

```bash
# Configuration: 
# SCALE=true
# SERVER_REGIONS=IL:gxy1,gxy2;US:gxy3,gxy4

# 1. First user from Israel requests room 2171
curl -X POST http://localhost:8081/v2/room_server \
  -H "Content-Type: application/json" \
  -d '{"room": 2171, "geo": {"country_code": "IL"}}'

# Response: {"janus": "gxy1"}  (regional server for IL)

# 2. Second user from US requests same room 2171
curl -X POST http://localhost:8081/v2/room_server \
  -H "Content-Type: application/json" \
  -d '{"room": 2171, "geo": {"country_code": "US"}}'

# Response: {"janus": "gxy1"}  (same server - sticky routing ignores country!)

# 3. After all users leave, assignment is removed
# 4. New first user from US requests room 2171
curl -X POST http://localhost:8081/v2/room_server \
  -H "Content-Type: application/json" \
  -d '{"room": 2171, "geo": {"country_code": "US"}}'

# Response: {"janus": "gxy3"}  (now gets US regional server)
```
