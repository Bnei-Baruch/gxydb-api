# MQTT Setup Guide

## Quick Start

### 1. Check Current Configuration

Check if MQTT is configured in your environment:

```bash
grep MQTT /etc/systemd/system/gxydb-api.service
# or
env | grep MQTT
```

### 2. Add MQTT Configuration

Add environment variables to gxydb-api configuration file:

```bash
# MQTT broker
MQTT_BROKER_URL=mqtt://mqtt.example.com:1883

# Client ID (unique for each instance)
MQTT_CLIENT_ID=gxydb-api-prod

# Admin secret for Janus (should already be configured)
GATEWAY_PLUGIN_ADMIN_KEY=your_admin_secret_here

# List of Janus servers to monitor (also used for room balancing)
AVAILABLE_JANUS_SERVERS=gxy1,gxy2,gxy3,gxy4,gxy5
```

**Important:** The `AVAILABLE_JANUS_SERVERS` list defines which servers will be:
- Monitored via MQTT (session metrics)
- Used for room balancing (in Scale Mode)

Only servers in this list will appear in Prometheus metrics.

### 3. Rebuild and Restart

```bash
cd /Users/amnonbb/go/src/github.com/Bnei-Baruch/gxydb-api
go build -o gxydb-api
sudo systemctl restart gxydb-api
```

### 4. Verify Operation

#### Check Logs
```bash
# Check MQTT connection
journalctl -u gxydb-api -n 100 | grep "Subscribed to gateway"

# Expected output:
# "Subscribed to gateway status" topic="janus/+/status"
# "Subscribed to gateway admin responses" topic="janus/+/from-janus-admin"

# Check status updates
journalctl -u gxydb-api -f | grep "Gateway"

# Expected output:
# "Gateway status updated" server="gxy1" online=true
# "Gateway sessions updated" server="gxy1" sessions=42
```

#### Check Prometheus Metrics
```bash
curl http://localhost:8081/metrics | grep gxydb_gateway_sessions

# Expected output:
# gxydb_gateway_sessions{gateway="gxy1",type="rooms"} 42
# gxydb_gateway_sessions{gateway="gxy2",type="rooms"} 15
```

## Janus Server Configuration

Each Janus server must publish data to MQTT. Example configuration (as in strdb):

### 1. Install MQTT Plugin for Janus

Check that Janus is built with MQTT transport:
```bash
janus --version | grep mqtt
```

### 2. Configure janus.transport.mqtt.jcfg

```json
{
    "general": {
        "enabled": true,
        "url": "tcp://mqtt.example.com:1883",
        "client_id": "gxy1",
        "username": "janus",
        "password": "your_mqtt_password",
        "subscribe_topic": "janus/gxy1/to-janus-admin",
        "publish_topic": "janus/gxy1/from-janus-admin"
    }
}
```

### 3. Status Publishing

Configure script to publish server status:

```bash
#!/bin/bash
# /usr/local/bin/janus-status.sh

MQTT_BROKER="mqtt.example.com"
MQTT_TOPIC="janus/gxy1/status"
SERVER_NAME="gxy1"

while true; do
    # Check if Janus is running
    if systemctl is-active --quiet janus; then
        mosquitto_pub -h $MQTT_BROKER -t $MQTT_TOPIC -m '{"online":true}' -r
    else
        mosquitto_pub -h $MQTT_BROKER -t $MQTT_TOPIC -m '{"online":false}' -r
    fi
    sleep 5
done
```

Create systemd service for this script:
```ini
# /etc/systemd/system/janus-status.service
[Unit]
Description=Janus Status Publisher
After=janus.service

[Service]
Type=simple
ExecStart=/usr/local/bin/janus-status.sh
Restart=always

[Install]
WantedBy=multi-user.target
```

Start it:
```bash
sudo systemctl enable janus-status
sudo systemctl start janus-status
```

## Verify MQTT Setup on Janus Side

### 1. Check Topic Subscriptions

```bash
# On MQTT broker
mosquitto_sub -h mqtt.example.com -t 'janus/+/status' -v

# Expected output (every 5 seconds):
# janus/gxy1/status {"online":true}
# janus/gxy2/status {"online":true}
```

### 2. Check Admin Requests

```bash
# On MQTT broker
mosquitto_sub -h mqtt.example.com -t 'janus/+/to-janus-admin' -v

# Expected output (every 10 seconds from gxydb-api):
# janus/gxy1/to-janus-admin {"janus":"list_sessions","transaction":"transaction","admin_secret":"..."}
```

### 3. Check Admin Responses

```bash
# On MQTT broker
mosquitto_sub -h mqtt.example.com -t 'janus/+/from-janus-admin' -v

# Expected output (Janus response to admin request):
# janus/gxy1/from-janus-admin {"janus":"success","transaction":"transaction","sessions":[...]}
```

## Troubleshooting

### Issue: "Subscribed to gateway" not appearing in logs

**Cause:** MQTT_BROKER_URL not configured or incorrect.

**Solution:**
```bash
# Check configuration
grep MQTT_BROKER_URL /etc/systemd/system/gxydb-api.service

# Check MQTT broker availability
telnet mqtt.example.com 1883
```

### Issue: "Gateway status updated" not appearing

**Cause:** Janus servers not publishing status to MQTT.

**Solution:**
```bash
# Check if status is being published
mosquitto_sub -h mqtt.example.com -t 'janus/+/status' -v -W 10

# If empty - check janus-status.service on Janus servers
systemctl status janus-status
journalctl -u janus-status -n 50
```

### Issue: "Gateway sessions updated" not appearing

**Cause:** Janus not responding to admin requests via MQTT.

**Solution:**
```bash
# Check MQTT transport in Janus
grep mqtt /opt/janus/etc/janus/janus.transport.mqtt.jcfg

# Check Janus logs
journalctl -u janus -n 100 | grep mqtt

# Manually send admin request
mosquitto_pub -h mqtt.example.com -t 'janus/gxy1/to-janus-admin' \
  -m '{"janus":"list_sessions","transaction":"test","admin_secret":"your_secret"}'

# Check response
mosquitto_sub -h mqtt.example.com -t 'janus/gxy1/from-janus-admin' -v -W 5
```

### Issue: Metrics show 0 for all servers

**Cause:** Either MQTT disabled or servers not responding.

**Solution:**
```bash
# Check gxydb-api logs
journalctl -u gxydb-api | grep -E "Gateway|MQTT"

# Check if MQTT is enabled
grep "MQTT disabled" /var/log/gxydb-api.log

# If MQTT disabled - add MQTT_BROKER_URL and restart
```

## Rollback to HTTP Admin API (if needed)

If MQTT doesn't work, you can temporarily return to HTTP Admin API:

```bash
# Remove MQTT_BROKER_URL from configuration
sudo vim /etc/systemd/system/gxydb-api.service
# Comment out: # Environment="MQTT_BROKER_URL=..."

# Restart
sudo systemctl daemon-reload
sudo systemctl restart gxydb-api

# Check logs
journalctl -u gxydb-api -n 50 | grep "HTTP Admin API"
```

**Note:** HTTP Admin API fallback works but shows warnings in logs. This is a temporary solution.

## Production Monitoring

### 1. Prometheus Alerts

```yaml
# prometheus-alerts.yml
groups:
  - name: gxydb-gateway
    rules:
      # Alert if no session data
      - alert: GatewayNoSessionsData
        expr: absent(gxydb_gateway_sessions)
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "No gateway sessions data"
          description: "gxydb-api is not collecting gateway sessions metrics"

      # Alert if gateway offline
      - alert: GatewayOffline
        expr: gxydb_gateway_sessions == 0
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "Gateway {{ $labels.gateway }} appears offline"
          description: "Gateway {{ $labels.gateway }} has 0 sessions for 10+ minutes"
```

### 2. Grafana Dashboard

Example queries for Grafana:

```promql
# Total sessions across all gateways
sum(gxydb_gateway_sessions)

# Sessions by server
gxydb_gateway_sessions{type="rooms"}

# Average load over 5 minutes
avg_over_time(gxydb_gateway_sessions[5m])
```

## Additional Information

- Detailed documentation: [MQTT_GATEWAY_MONITORING.md](./MQTT_GATEWAY_MONITORING.md)
- Changes summary: [MQTT_IMPLEMENTATION_SUMMARY.md](./MQTT_IMPLEMENTATION_SUMMARY.md)
- Failover documentation: [FAILOVER.md](./FAILOVER.md)
- Reference project: `/Users/amnonbb/go/src/github.com/Bnei-Baruch/strdb`
