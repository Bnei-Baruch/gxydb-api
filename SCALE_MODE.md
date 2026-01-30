# Scale Mode Configuration

Quick reference for switching between legacy and scale modes.

## TL;DR

```bash
# Legacy mode (default) - use room's default gateway
SCALE=false

# Scale mode - dynamic load balancing
SCALE=true
```

## Configuration Files

### .env for Legacy Mode (Default)

```bash
# No SCALE variable needed - defaults to false
# OR explicitly:
SCALE=false
```

### .env for Scale Mode

```bash
# Enable scale mode
SCALE=true

# Configure servers
AVAILABLE_JANUS_SERVERS=gxy1,gxy2,gxy3,gxy4,gxy5,gxy6,gxy7,gxy8,gxy9,gxy10,gxy11,gxy12
MAX_SERVER_CAPACITY=400
AVG_ROOM_OCCUPANCY=10

# Optional: Regional routing
SERVER_REGIONS=IL:gxy1,gxy2,gxy3;US:gxy4,gxy5,gxy6;RU:gxy7,gxy8,gxy9
```

## Behavior Comparison

| Feature | SCALE=false (Legacy) | SCALE=true (Scale) |
|---------|---------------------|-------------------|
| Server selection | Room's default_gateway | Load balanced |
| Sticky routing | Always same server | First user sets server |
| Regional routing | Not supported | Supported via geo.country_code |
| Database writes | None | room_server_assignments table |
| Auto-cleanup | N/A | Yes, on session end |
| Monitoring | Simple | Detailed (load, region, etc) |

## Migration Steps

### Safe Rollout

```bash
# 1. Deploy with legacy mode
SCALE=false
# → All clients use room's default gateway
# → Zero risk, existing behavior

# 2. Test scale mode on staging
SCALE=true
# → Verify load balancing works
# → Test regional routing
# → Monitor assignment cleanup

# 3. Enable on production
SCALE=true
# → Switch at low-traffic time
# → Monitor logs for "Assigned room to server"
# → Watch server load distribution

# 4. Rollback if needed (instant)
SCALE=false
# → Revert to legacy behavior immediately
# → No data migration needed
```

### Testing Checklist

- [ ] SCALE=false: All rooms return their default gateway
- [ ] SCALE=false: Verify all rooms have valid `default_gateway_id`
- [ ] SCALE=false: Test error handling for missing gateway
- [ ] SCALE=true without regions: Load balancing works
- [ ] SCALE=true with regions: Regional routing works
- [ ] Sticky routing: Second user gets same server
- [ ] Auto-cleanup: Assignments removed when room empty
- [ ] Monitoring: Check logs for assignment events

## Monitoring

### Legacy Mode (SCALE=false)

Watch for gateway configuration errors:

```
ERROR: Gateway not found for room in legacy mode (SCALE=false). 
       Ensure all rooms have valid default_gateway_id.
       room_id=2171 default_gateway_id=999 room_gateway_uid=2171
```

**Action required:** Fix room configuration in database before clients can use that room.

### Scale Mode (SCALE=true)

Watch for these log messages:

```
INFO: Assigned room to server room_id=2171 gateway_name=gxy5 country_code=IL regional_match=true
INFO: Cleaned up room assignment immediately room_id=2171
INFO: Cleaned inactive room server assignments cleaned_assignments=5
```

### Metrics to Monitor

- Server load distribution (COUNT rooms per gateway × AVG_ROOM_OCCUPANCY)
- Assignment creation rate
- Assignment cleanup rate
- Regional match rate (if using SERVER_REGIONS)
- Rooms per server (should not exceed MAX_SERVER_CAPACITY / AVG_ROOM_OCCUPANCY)

## Troubleshooting

### Issue: Error "gateway not found for room" in legacy mode

**Cause:** Room's `default_gateway_id` is missing or invalid in SCALE=false mode.

**Solution:**
```sql
-- Find rooms without valid gateway
SELECT r.id, r.name, r.default_gateway_id 
FROM rooms r 
LEFT JOIN gateways g ON r.default_gateway_id = g.id 
WHERE g.id IS NULL;

-- Fix: assign valid gateway to room
UPDATE rooms SET default_gateway_id = (SELECT id FROM gateways WHERE name = 'gxy1' LIMIT 1) WHERE id = 2171;
```

### Issue: All requests go to same server in scale mode

**Check:**
1. Is SCALE=true set?
2. Are AVAILABLE_JANUS_SERVERS configured?
3. Check logs for "Assigned room to server"

### Issue: Regional routing not working

**Check:**
1. Is SERVER_REGIONS formatted correctly?
2. Is geo.country_code sent in request?
3. Check logs for "regional_match=true/false"

### Issue: Need to rollback

**Solution:**
```bash
# Immediate rollback to legacy mode
SCALE=false
# Restart service
```

No data cleanup needed - room_server_assignments table is simply ignored.

## Performance Notes

### Legacy Mode (SCALE=false)
- 🟢 Zero overhead
- 🟢 No database writes
- 🟢 Fast cache lookup only
- ⚠️  Requires valid `default_gateway_id` for all rooms (strict validation)

### Scale Mode (SCALE=true)
- 🟡 DB read on first request (check existing assignment)
- 🟡 DB write on first user (create assignment)
- 🟡 DB write on session close (cleanup if last user)
- 🟢 Cache lookup for gateway load calculation

Expected overhead: < 5ms per request
