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
  "room": 2171
}
```

**Response:**
```json
{
  "janus": "gxy5"
}
```

**Response codes:**
- `200 OK` - server successfully retrieved
- `404 Not Found` - room not found
- `500 Internal Server Error` - server error

## Configuration

### Environment Variable

```bash
AVAILABLE_JANUS_SERVERS=gxy1,gxy2,gxy3,gxy4,gxy5,gxy6,gxy7,gxy8,gxy9,gxy10,gxy11,gxy12
```

By default, all servers gxy1-12 are used.

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

- **Maximum server capacity**: 400 users (hardcoded)
- **Average room capacity**: 10 users (for calculations)
- **Server priorities**: Not implemented in current version

## Future Improvements

1. **Regional load balancing**: 
   - Use `region` field in table
   - Server assignment based on user region

2. **Server priorities**:
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

## Usage Example

```bash
# 1. Client wants to connect to room 2171
curl -X POST http://localhost:8081/v2/room_server \
  -H "Content-Type: application/json" \
  -d '{"room": 2171}'

# Response: {"janus": "gxy5"}

# 2. Client connects to server gxy5
# 3. All subsequent users for room 2171 will get the same server "gxy5"

# 4. After all users leave the room, the assignment will be automatically removed
# 5. Next request for room 2171 may get a different server depending on current load
```
