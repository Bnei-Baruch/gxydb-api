# Difference between /groups and /rooms

## GET /groups - Static Assignments

**Returns:** Factory settings from `rooms` table

```json
{
  "rooms": [
    {
      "room": 2171,
      "janus": "gxy1",        ← default_gateway from rooms.default_gateway_id
      "description": "Room 1",
      "num_users": 5
    }
  ]
}
```

**Features:**
- ✅ Fast (cached)
- ✅ Shows "factory" settings
- ✅ Does NOT depend on SCALE mode
- ✅ Used for listing available rooms

---

## GET /rooms - Dynamic Assignments

**Returns:** Real state with dynamic assignments

```json
[
  {
    "room": 2171,
    "janus": "gxy5",         ← from room_server_assignments (if SCALE=true)
                             ← or default_gateway (if SCALE=false or no assignment)
    "description": "Room 1",
    "users": [
      {"id": "user1", "display": "..."},
      {"id": "user2", "display": "..."}
    ]
  }
]
```

**Features:**
- ⚠️ Slower (DB JOIN)
- ✅ Shows real assignments
- ✅ Depends on SCALE mode
- ✅ Full user information
- ✅ Used for monitoring active rooms

---

## /rooms Logic:

```
SCALE=false (Legacy Mode):
  → janus = room.default_gateway
  
SCALE=true (Scale Mode):
  1. Check room_server_assignments for room_id
  2. If found → janus = gateway_name (dynamic)
  3. If not found → janus = default_gateway (static)
```

---

## Usage Examples:

### Show all rooms list (static):
```bash
GET /groups?with_num_users=true
```

### Show active rooms with real assignments:
```bash
GET /rooms
```

### Check which server is assigned to a room:
```bash
# Static assignment:
GET /groups → "janus": "gxy1"

# Real assignment (if there are users):
GET /rooms → "janus": "gxy5"
```

---

## Scenario:

```
1. Room created with default_gateway_id = gxy1
   GET /groups → "janus": "gxy1" ✅
   GET /rooms → no active users → not shown

2. User requests server (SCALE=true):
   POST /v2/room_server → gxy5 assigned
   room_server_assignments: room_id=100, gateway_name=gxy5

3. User enters:
   GET /groups → "janus": "gxy1" (unchanged)
   GET /rooms → "janus": "gxy5" ✅ (shows real assignment)

4. Users leave:
   room_server_assignments deleted
   GET /groups → "janus": "gxy1" (unchanged)
   GET /rooms → no active users → not shown

5. New user enters:
   POST /v2/room_server → gxy2 assigned (can be different!)
   GET /rooms → "janus": "gxy2" (new assignment)
```

---

## Important:

- **GET /groups** always shows static settings
- **GET /rooms** only shows rooms with active users
- **GET /rooms** in Scale Mode shows real dynamic assignments
- Dynamic assignments are temporary and deleted when all users leave
