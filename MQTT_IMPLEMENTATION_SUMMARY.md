# MQTT Implementation Summary

## Task

Replace Janus gateway server polling via HTTP Admin API with MQTT transport, based on the `strdb` project example.

## Modified Files

### 1. `api/mqtt.go`
**Changes:**
- Added `GatewayStatus` struct for storing server state
- Added `gatewayStatuses map[string]*GatewayStatus` field to `MQTTListener`
- Added `adminSecret` field for admin requests
- Updated `NewMQTTListener` constructor - added `adminSecret` parameter
- Added subscription to topics:
  - `janus/+/status` - server status (online/offline)
  - `janus/+/from-janus-admin` - admin request responses
- Added handlers:
  - `HandleGatewayStatus()` - processes status messages
  - `HandleGatewayAdminResponse()` - processes admin responses
- Added `startPeriodicAdminMessages()` method - sends `list_sessions` every 10 seconds
- Added `SendAdminMessage()` method - sends admin messages
- Added `GetGatewayStatuses()` method - returns statuses for external use

**New types:**
```go
type GatewayStatus struct {
    Name     string
    Online   bool
    Sessions int
    LastSeen time.Time
    mu       sync.RWMutex
}
```

### 2. `api/app.go`
**Changes:**
- Updated `NewMQTTListener` call - added `common.Config.GatewayPluginAdminKey` parameter
- Updated `NewPeriodicCollector` call - added `a.mqttListener` parameter

### 3. `instrumentation/periodic.go`
**Changes:**
- Removed `janus_admin` and `models` imports (no longer needed)
- Added `GatewayStatusProvider` interface:
  ```go
  type GatewayStatusProvider interface {
      GetGatewayStatuses() map[string]*common.GatewayStatusInfo
  }
  ```
- Updated `NewPeriodicCollector` constructor - added `gatewayStatusProvider` parameter
- Completely rewrote `collectGatewaySessions()` method:
  - Uses server list from `common.Config.AvailableJanusServers` instead of DB
  - If MQTT enabled (gatewayStatusProvider != nil) - uses data from MQTT
  - If MQTT disabled - returns error (HTTP Admin API deprecated)
  - Added logging for offline and stale servers
  - All servers get `type="rooms"` label in metrics

### 4. `instrumentation/periodic_test.go`
**Changes:**
- Updated `NewPeriodicCollector(s.DB, nil)` call - added `nil` for gatewayStatusProvider

### 5. `common/consts.go`
**Changes:**
- Added `GatewayStatusInfo` struct for passing data between packages:
  ```go
  type GatewayStatusInfo struct {
      Name     string
      Online   bool
      Sessions int
      LastSeen time.Time
  }
  ```

## New Files

### 1. `MQTT_GATEWAY_MONITORING.md`
Detailed documentation about:
- MQTT monitoring architecture
- MQTT topics and message formats
- Configuration
- Monitoring and logs
- Migration from HTTP Admin API

### 2. `MQTT_IMPLEMENTATION_SUMMARY.md`
This file - brief summary of changes.

## Workflow

### 1. Initialization (on gxydb-api startup)
1. `initMQTT()` creates `MQTTListener` with adminSecret
2. `MQTTListener.Start()` connects to MQTT broker
3. Subscribes to topics `janus/+/status` and `janus/+/from-janus-admin`
4. Starts goroutine `startPeriodicAdminMessages()`
5. `initInstrumentation()` creates `PeriodicCollector` with reference to `MQTTListener`

### 2. Server Status Monitoring
1. Janus server publishes `{"online": true}` to `janus/{server}/status`
2. `HandleGatewayStatus()` receives message
3. Updates `gatewayStatuses[server].Online = true`

### 3. Session Metrics Collection
1. Every 10 seconds `startPeriodicAdminMessages()`:
   - Iterates through all servers in `gatewayStatuses`
   - If `Online == true` → sends `list_sessions` to `janus/{server}/to-janus-admin`
2. Janus server responds with `{"janus": "success", "sessions": [...]}` to `janus/{server}/from-janus-admin`
3. `HandleGatewayAdminResponse()` receives response
4. Updates `gatewayStatuses[server].Sessions = len(response.Sessions)`

### 4. Prometheus Metrics Publishing
1. Every second `PeriodicCollector.collectGatewaySessions()`:
   - Iterates through servers from `AVAILABLE_JANUS_SERVERS` config
   - Calls `mqttListener.GetGatewayStatuses()` to get status for each server
   - Updates `Stats.GatewaySessionsGauge` for each server
   - Logs warnings for offline or stale servers

## Benefits

1. **Performance**: Fewer HTTP requests, fewer timeouts
2. **Real-time**: Instant server status updates
3. **Scalability**: Easy to add new servers (just update `AVAILABLE_JANUS_SERVERS`)
4. **Reliability**: MQTT auto-reconnect, graceful degradation
5. **Consistency**: Same server list used for both monitoring and room balancing

## Backward Compatibility

- If `MQTT_BROKER_URL` not configured - system works without gateway session monitoring
- All other gxydb-api functions work as before

## Dependencies

- Library: `github.com/eclipse/paho.mqtt.golang` (already in project)
- MQTT broker (e.g., Mosquitto, RabbitMQ, EMQ)

## Testing

### Unit tests
```bash
go test ./instrumentation/... -v
```

### Build
```bash
go build -o gxydb-api
```

### Check MQTT connection
```bash
# Should see in logs:
grep "Subscribed to gateway" /var/log/gxydb-api.log
```

### Check metrics
```bash
curl http://localhost:8081/metrics | grep gxydb_gateway_sessions
```

## Next Steps

1. Configure MQTT broker in test environment
2. Configure Janus servers to publish status and admin responses
3. Test in staging environment
4. Deploy to production
5. Monitor logs and metrics
6. After stabilization - remove HTTP Admin API fallback code
