# Janus Gateway Failover

## Overview

Automatic failover mechanism that migrates room assignments from failed primary servers to standby failover servers.

## Configuration

```bash
# Primary servers for normal operation and room balancing
AVAILABLE_JANUS_SERVERS=gxy1,gxy2,gxy3,gxy4,gxy5

# Standby failover servers (not used for normal balancing)
FAILOVER_JANUS_SERVERS=gxy12,gxy13

# Wait time before triggering failover (default: 5s)
FAILOVER_WAIT_TIME=5s
```

## Behavior

### When Primary Server Goes Offline

1. **Detection**: MQTTListener receives `online: false` status message
2. **Wait Period**: Starts timer for `FAILOVER_WAIT_TIME` (default 5 seconds)
3. **Failover Trigger**: After wait period expires:
   - Selects next available failover server
   - Migrates all room assignments from failed server to failover
   - Records failover mapping
   - Updates metrics

### When Failed Server Recovers

1. **Detection**: MQTTListener receives `online: true` status message
2. **Graceful Recovery** (natural migration):
   - Cancels any pending failover timer
   - New room assignments go to recovered server
   - Existing assignments stay on failover
   - Releases failover server back to available pool
   - Assignments naturally migrate as users leave rooms

### When Failover Servers Exhausted

If all failover servers are in use and another primary fails:

1. **Emergency Distribution**: Assignments distributed to alive primary servers
2. **Capacity Override**: Ignores `MAX_SERVER_CAPACITY` limits
3. **Round-Robin**: Distributes assignments evenly among alive servers
4. **Metrics**: Logs CRITICAL and increments `failover_events_total{event="no_failover"}`

## Server Types

### Primary Servers (`AVAILABLE_JANUS_SERVERS`)
- Used for normal room balancing
- Participate in dynamic assignments
- Monitored for failover
- Can receive emergency assignments

### Failover Servers (`FAILOVER_JANUS_SERVERS`)
- Reserved for emergencies only
- Do NOT participate in normal room balancing
- Only used when primary server fails
- Monitored alongside primary servers

## Prometheus Metrics

### Failover Events Counter
```
galaxy_failover_events_total{event, failed_server, failover_server}
```

Events:
- `triggered` - Failover process started
- `completed` - Failover successfully completed
- `failed` - Failover process failed
- `no_failover` - No failover servers available
- `distributed` - Emergency distribution completed
- `recovered` - Failed server came back online

### Active Failover Gauge
```
galaxy_failover_active{failed_server, failover_server}
```

Values:
- `1` - Failover mapping active
- `0` - Failover mapping inactive (server recovered)

### Gateway Sessions Gauge
```
galaxy_gateways_sessions{name, type}
```

Includes both AVAILABLE and FAILOVER servers.

## Failure Scenarios

### Scenario 1: Single Server Failure
```
Initial: gxy1 (online), gxy2 (online), gxy12 (standby)
↓
gxy1 fails (online: false)
↓
Wait 5 seconds
↓
Migrate: gxy1 → gxy12
  - 40 room assignments moved
  - gxy12 marked as occupied
↓
State: gxy1 (failed, assignments on gxy12), gxy2 (online), gxy12 (active failover)
```

### Scenario 2: Server Recovery
```
gxy1 recovers (online: true)
↓
- Cancel failover timer (if still pending)
- New assignments → gxy1
- Existing assignments stay on gxy12
- Release gxy12 to available pool
↓
As users leave rooms:
  - Assignments naturally expire
  - gxy12 becomes empty
  - gxy12 ready for next failover
```

### Scenario 3: Multiple Failures (Exhaustion)
```
Initial: gxy1, gxy2, gxy3 (online), gxy12, gxy13 (standby)
↓
gxy1 fails → gxy12 takes over
gxy2 fails → gxy13 takes over
gxy3 fails → NO FAILOVER AVAILABLE
↓
Emergency: Distribute gxy3's assignments to gxy1 and gxy2
  - Round-robin distribution
  - Ignore capacity limits
  - Log CRITICAL
```

## Logs

### Failover Triggered
```json
{
  "level": "warn",
  "server": "gxy1",
  "wait_time": "5s",
  "message": "Primary server offline - failover timer started"
}
```

### Failover Completed
```json
{
  "level": "info",
  "failed_server": "gxy1",
  "failover_server": "gxy12",
  "count": 40,
  "message": "Failover completed successfully"
}
```

### Server Recovered
```json
{
  "level": "info",
  "server": "gxy1",
  "failover_server": "gxy12",
  "message": "Failed server recovered - new assignments will go to recovered server"
}
```

### No Failover Available
```json
{
  "level": "error",
  "server": "gxy3",
  "message": "No failover servers available - will distribute to alive primary servers"
}
```

## Monitoring

### Grafana Alerts

```yaml
# Alert when failover triggered
- alert: FailoverTriggered
  expr: increase(galaxy_failover_events_total{event="triggered"}[5m]) > 0
  labels:
    severity: warning
  annotations:
    summary: "Failover triggered for {{ $labels.failed_server }}"

# Alert when failover exhausted
- alert: FailoverExhausted
  expr: increase(galaxy_failover_events_total{event="no_failover"}[5m]) > 0
  labels:
    severity: critical
  annotations:
    summary: "No failover servers available"

# Alert when server recovered
- alert: ServerRecovered
  expr: increase(galaxy_failover_events_total{event="recovered"}[5m]) > 0
  labels:
    severity: info
  annotations:
    summary: "Server {{ $labels.failed_server }} recovered"
```

### Grafana Dashboard

```promql
# Active failovers
sum(galaxy_failover_active)

# Failover events by type
rate(galaxy_failover_events_total[5m])

# Failed servers list
galaxy_failover_active == 1
```

## Testing

### Test Failover Trigger

```bash
# 1. Publish offline status for primary server
mosquitto_pub -h mqtt.example.com -t 'janus/gxy1/status' -m '{"online":false}' -r

# 2. Wait 5 seconds and check logs
journalctl -u gxydb-api -f | grep -E "failover|gxy1"

# 3. Check metrics
curl http://localhost:8081/metrics | grep failover

# 4. Verify assignment migration
psql galaxy -c "SELECT gateway_name, COUNT(*) FROM room_server_assignments GROUP BY gateway_name;"
```

### Test Recovery

```bash
# 1. Publish online status
mosquitto_pub -h mqtt.example.com -t 'janus/gxy1/status' -m '{"online":true}' -r

# 2. Check logs for recovery
journalctl -u gxydb-api -n 50 | grep recovered

# 3. Verify new assignments go to recovered server
curl -X POST http://localhost:8081/v2/room_server -d '{"room": 999}'
```

## Best Practices

1. **Failover Count**: Configure at least 2 failover servers for redundancy
2. **Wait Time**: 5-10 seconds is recommended (balance between false positives and recovery speed)
3. **Monitoring**: Set up alerts for all failover events
4. **Testing**: Regularly test failover in staging environment
5. **Capacity Planning**: Ensure failover servers have same capacity as primary

## Limitations

1. **No Automatic Failback**: Assignments don't automatically move back when primary recovers
2. **Retained Messages**: Janus servers must publish status with `retained` flag for instant detection
3. **Capacity Ignored**: Emergency distribution ignores capacity limits (by design)
4. **Network Split**: Cannot distinguish between server failure and network partition

## Troubleshooting

### Failover Not Triggering

**Check:**
- MQTT connection: `journalctl -u gxydb-api | grep "Subscribed to gateway"`
- Status messages: `mosquitto_sub -h mqtt.example.com -t 'janus/+/status' -v`
- Failover config: `grep FAILOVER /etc/systemd/system/gxydb-api.service`

### Assignments Not Migrating

**Check:**
- Room assignment manager initialized: `grep "roomServerAssignmentManager" /var/log/gxydb-api.log`
- Database connectivity
- Error logs: `journalctl -u gxydb-api | grep "migrate assignments"`

### Failover Server Not Released

**Check:**
- Server actually recovered: `mosquitto_sub -h mqtt.example.com -t 'janus/gxy1/status'`
- Mapping status: Check `galaxy_failover_active` metric
- Logs: `journalctl -u gxydb-api | grep recovered`

## Related Documentation

- [MQTT Gateway Monitoring](./MQTT_GATEWAY_MONITORING.md)
- [Room Server Balancing](./ROOM_SERVER_BALANCING.md)
- [Scale Mode](./SCALE_MODE.md)
