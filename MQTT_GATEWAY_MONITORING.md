# MQTT Gateway Monitoring

## Overview

Janus gateway server monitoring system via MQTT transport instead of HTTP Admin API.

## Advantages of MQTT Approach

1. **Lower load**: Single MQTT broker instead of multiple HTTP requests to each server
2. **Real-time status**: Status topic subscription allows instant notification of server issues
3. **Fewer timeouts**: No need to wait for HTTP timeouts on unavailable servers
4. **Centralized management**: All servers publish data to a single MQTT broker

## Architecture

### MQTT Topics

1. **Status topics** (`janus/{server}/status`):
   - Server publishes its status: `{"online": true}` or `{"online": false}`
   - gxydb-api subscribes to `janus/+/status`
   - Updates `Online` field in server status

2. **Admin request topics** (`janus/{server}/to-janus-admin`):
   - gxydb-api sends `list_sessions` request every 10 seconds
   - Only to servers with `Online = true` status

3. **Admin response topics** (`janus/{server}/from-janus-admin`):
   - Server responds with session list: `{"janus": "success", "sessions": [...]}`
   - gxydb-api subscribes to `janus/+/from-janus-admin`
   - Updates `Sessions` field in server status

### Components

#### MQTTListener (api/mqtt.go)

- Connects to MQTT broker
- Subscribes to status and admin topics
- Stores all gateway statuses in `gatewayStatuses map[string]*GatewayStatus`
- Periodically sends `list_sessions` requests

Key methods:
- `HandleGatewayStatus()` - processes status messages
- `HandleGatewayAdminResponse()` - processes admin responses
- `SendAdminMessage()` - sends admin requests
- `GetGatewayStatuses()` - returns statuses for PeriodicCollector

#### PeriodicCollector (instrumentation/periodic.go)

- Collects metrics for Prometheus
- Gets data from `MQTTListener.GetGatewayStatuses()`
- Uses server list from `AVAILABLE_JANUS_SERVERS` config
- If MQTT disabled - fallback to HTTP Admin API (deprecated)

### Gateway Status

```go
type GatewayStatus struct {
    Name     string    // Server name (gxy1, gxy2, ...)
    Online   bool      // Is server online (from status topic)
    Sessions int       // Number of sessions (from admin response)
    LastSeen time.Time // Last update time
}
```

## Configuration

### Environment Variables

```bash
# MQTT broker
MQTT_BROKER_URL=mqtt://mqtt.example.com:1883

# Client ID for gxydb-api
MQTT_CLIENT_ID=gxydb-api

# Admin secret for Janus (for list_sessions requests)
GATEWAY_PLUGIN_ADMIN_KEY=your_admin_secret

# Primary servers for room balancing (also monitored)
AVAILABLE_JANUS_SERVERS=gxy1,gxy2,gxy3,gxy4,gxy5

# Optional: Failover servers for high availability
FAILOVER_JANUS_SERVERS=gxy12,gxy13
FAILOVER_WAIT_TIME=5s
```

**Note:** Both `AVAILABLE_JANUS_SERVERS` and `FAILOVER_JANUS_SERVERS` are monitored via MQTT. For failover behavior, see [FAILOVER.md](./FAILOVER.md).

### Configuration on Janus Servers

Each Janus server must:

1. Publish status to `janus/{server_name}/status`
2. Subscribe to `janus/{server_name}/to-janus-admin`
3. Publish responses to `janus/{server_name}/from-janus-admin`

Example configuration (similar to strdb): see `/Users/amnonbb/go/src/github.com/Bnei-Baruch/strdb`

## Monitoring

### Prometheus Metrics

- `galaxy_gateways_sessions{name="gxy1", type="rooms"}` - number of sessions on server
- `galaxy_failover_events_total{event, failed_server, failover_server}` - failover events counter
- `galaxy_failover_active{failed_server, failover_server}` - active failover mappings

### Logs

```bash
# Topic subscription
"Subscribed to gateway status" topic="janus/+/status"
"Subscribed to gateway admin responses" topic="janus/+/from-janus-admin"

# Status updates
"Gateway status updated" server="gxy1" online=true

# Session updates
"Gateway sessions updated" server="gxy1" sessions=42

# Warnings
"Gateway is offline" gateway="gxy1"
"Gateway status is stale" gateway="gxy1" since_last_seen=45s
```

## Differences from HTTP Admin API

### Old Approach (HTTP):
- Send HTTP request to each server every second
- 900ms timeout for entire server pool
- No server status information until timeout

### New Approach (MQTT):
- Server publishes status itself (real-time)
- `list_sessions` requests sent only to online servers (every 10 seconds)
- No timeouts - MQTT broker handles delivery

## Failover and Reliability

1. **MQTT auto-reconnect**: Client automatically reconnects on connection loss
2. **Status tracking**: Warning logged if server doesn't update status for > 30 seconds
3. **Graceful fallback**: If MQTT disabled (no MQTT_BROKER_URL) - system works without monitoring

## Migration

### Steps to Enable MQTT Monitoring:

1. Configure MQTT broker
2. Configure Janus servers to publish status and admin responses
3. Add environment variables to gxydb-api
4. Restart gxydb-api
5. Check logs: `grep "Subscribed to gateway" /var/log/gxydb-api.log`
6. Check metrics: `curl http://localhost:8081/metrics | grep gxydb_gateway_sessions`

### Rollback (if something goes wrong):

- Remove `MQTT_BROKER_URL` from environment variables
- Restart gxydb-api
- System continues to work (without gateway session monitoring)

## High Availability

For automatic failover support, see [FAILOVER.md](./FAILOVER.md).

## Compatibility

- Implementation compatible with approach from `strdb` project
- Uses same topics and message formats
- Uses same library: `github.com/eclipse/paho.mqtt.golang`
