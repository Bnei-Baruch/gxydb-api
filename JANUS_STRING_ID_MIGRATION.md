# Janus String ID Migration Guide

## Overview
This migration converts room_id and rfid (user_id) from integers to strings to support Janus Gateway string IDs and Keycloak UUIDs.

## What Changes

### Database Schema Changes:

1. **sessions table:**
   - `room_id`: BIGINT → VARCHAR(255) (stores Janus room ID, e.g., "2171")
   - `gateway_feed` (RFID): BIGINT → VARCHAR(255) (stores user UUID from Keycloak)

2. **room_statistics table:**
   - `room_id`: BIGINT → VARCHAR(255)

3. **composites_rooms table:**
   - `room_id`: BIGINT → VARCHAR(255)

4. **rooms table:**
   - `gateway_uid`: INTEGER → VARCHAR(255) (can now store UUID room IDs)

5. **room_server_assignments table (new):**
   - `room_id`: VARCHAR(255) (stores Janus room ID for load balancing)

### Code Changes:
- API types updated to use `string` instead of `int` for room IDs
- Cache layer converts int → string automatically
- Load balancing uses Janus room IDs instead of internal DB IDs

## Migration Steps

### 1. Apply Database Migrations

```bash
# Set your database connection
export DATABASE_URL="postgres://user:password@localhost:5432/gxydb?sslmode=disable"

# Apply migrations (they will run in order)
migrate -path migrations -database "$DATABASE_URL" up

# Verify migrations applied
migrate -path migrations -database "$DATABASE_URL" version
```

**Migrations applied:**
1. `20240321_change_rfid_roomid_to_string` - sessions, room_statistics, composites_rooms
2. `20240322_janus_string_room_id` - rooms.gateway_uid
3. `20260119120000_room_server_assignments` - new table for load balancing

### 2. Regenerate SQLBoiler Models

```bash
# Install sqlboiler if not already installed
go install github.com/volatiletech/sqlboiler/v4@latest
go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest

# Regenerate models to match new schema
sqlboiler psql

# This will update models/ directory with correct types
```

### 3. Rebuild Application

```bash
# Build
go build -o gxydb-api .

# Run tests
go test ./...
```

### 4. Deploy

```bash
# Stop current service
systemctl stop gxydb-api

# Deploy new binary
cp gxydb-api /path/to/production/

# Start service
systemctl start gxydb-api
```

## Rollback

If you need to rollback:

```bash
# Rollback 3 migrations
migrate -path migrations -database "$DATABASE_URL" down 3

# Regenerate old models
sqlboiler psql

# Rebuild with old code
git checkout <previous-commit>
go build -o gxydb-api .
```

## Testing

### Test API Endpoint

```bash
# Get rooms list
curl http://localhost:8081/galaxy/groups | jq '.rooms[0]'

# Expected output:
# {
#   "room": "2171",    # ← now string!
#   "janus": "gxy5",
#   "description": "Africa"
# }

# Request server assignment (SCALE=true)
curl -X POST http://localhost:8081/v2/room_server \
  -H "Content-Type: application/json" \
  -d '{"room": "2171"}'

# Expected output:
# {"janus": "gxy1"}
```

### Verify Database

```sql
-- Check sessions table
SELECT room_id, gateway_feed FROM sessions LIMIT 5;
-- room_id should be string like "2171"
-- gateway_feed should be string UUID

-- Check rooms table  
SELECT gateway_uid FROM rooms LIMIT 5;
-- gateway_uid should be string like "2171"

-- Check room_server_assignments
SELECT * FROM room_server_assignments;
-- room_id should be string
```

## Important Notes

1. **Backward Compatibility:**
   - Old integer room IDs (2171) are automatically converted to strings ("2171")
   - Clients receiving room list will get strings
   - Clients sending the same strings back works seamlessly

2. **Future UUID Support:**
   - After this migration, you can add rooms with UUID gateway_uids
   - Example: `INSERT INTO rooms (gateway_uid, ...) VALUES ('room-uuid-123', ...)`
   - Janus will accept both "2171" and "room-uuid-123" formats

3. **Load Balancing:**
   - Uses Janus room IDs (gateway_uid) for tracking, not internal DB IDs
   - Sticky routing works correctly with string IDs

## Verification Checklist

- [ ] All 3 migrations applied successfully
- [ ] SQLBoiler models regenerated
- [ ] Application builds without errors
- [ ] Tests pass
- [ ] API returns rooms with string IDs
- [ ] Can create new sessions with string room_id
- [ ] Load balancing works (if SCALE=true)
- [ ] Existing data migrated correctly

## Support

If you encounter issues:
1. Check migration status: `migrate -path migrations -database "$DATABASE_URL" version`
2. Check logs: `journalctl -u gxydb-api -f`
3. Verify schema: `\d sessions` in psql

## Summary

This migration enables full Janus String ID support while maintaining backward compatibility with existing integer room IDs. All room IDs are now strings throughout the system, allowing future use of UUIDs for rooms and users.
