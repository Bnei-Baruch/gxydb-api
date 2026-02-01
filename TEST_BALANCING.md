# Room Balancing Test

## Test Configuration

```bash
SCALE=true
AVAILABLE_JANUS_SERVERS=gxy1,gxy2,gxy3
MAX_SERVER_CAPACITY=400
AVG_ROOM_OCCUPANCY=10
```

## Expected Behavior

**Formula:** Rooms per server = MAX_SERVER_CAPACITY / AVG_ROOM_OCCUPANCY = 400 / 10 = **40 rooms**

## Test: Create 120 Rooms

```bash
# Create 120 rooms sequentially
for i in {1..120}; do
  curl -X POST http://localhost:8081/v2/room_server \
    -H "Content-Type: application/json" \
    -d "{\"room\": $i}"
done
```

### Expected Result (Sequential Filling):

```
Rooms 1-40   → gxy1 (filled to limit: 400 users)
Rooms 41-80  → gxy2 (filled to limit: 400 users)
Rooms 81-120 → gxy3 (filled to limit: 400 users)

Strategy: fill servers SEQUENTIALLY, not evenly
```

## Check Distribution

```sql
-- Number of rooms on each server
SELECT gateway_name, COUNT(*) as rooms_count, COUNT(*) * 10 as estimated_load
FROM room_server_assignments
GROUP BY gateway_name
ORDER BY gateway_name;

-- Expected result:
-- gateway_name | rooms_count | estimated_load
-- gxy1         | 40          | 400
-- gxy2         | 40          | 400
-- gxy3         | 40          | 400
```

## Check Limit

```bash
# Create 121st room (all servers full)
curl -X POST http://localhost:8081/v2/room_server \
  -H "Content-Type: application/json" \
  -d '{"room": 121}'

# Result: gxy1 (fallback - selects least loaded)
```

## Important

- Balancing works **proactively** - no need to wait for users to enter
- Even if each room has only 1 real user, distribution will be even
- MAX_SERVER_CAPACITY and AVG_ROOM_OCCUPANCY **are used** for calculations
