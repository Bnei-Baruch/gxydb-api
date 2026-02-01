package api

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
	"github.com/edoshor/janus-go"
	"net/url"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	pkgerr "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Bnei-Baruch/gxydb-api/common"
)

// GatewayStatus holds the status of a Janus gateway
type GatewayStatus struct {
	Name     string
	Online   bool
	Sessions int
	LastSeen time.Time
	mu       sync.RWMutex
}

type MQTTListener struct {
	client                 mqtt.Client
	cache                  *AppCache
	serviceProtocolHandler ServiceProtocolHandler
	SessionManager         SessionManager
	gatewayStatuses        map[string]*GatewayStatus
	gatewayStatusesMu      sync.RWMutex
	periodicTicker         *time.Ticker
	adminSecret            string
}

func NewMQTTListener(cache *AppCache, sph ServiceProtocolHandler, sm SessionManager, adminSecret string) *MQTTListener {
	return &MQTTListener{
		cache:                  cache,
		serviceProtocolHandler: sph,
		SessionManager:         sm,
		gatewayStatuses:        make(map[string]*GatewayStatus),
		adminSecret:            adminSecret,
	}
}

func (l *MQTTListener) Start() error {
	// TODO: take log level from config
	// logging
	mqtt.DEBUG = NewPahoLogAdapter(zerolog.DebugLevel)
	mqtt.WARN = NewPahoLogAdapter(zerolog.WarnLevel)
	mqtt.CRITICAL = NewPahoLogAdapter(zerolog.ErrorLevel)
	mqtt.ERROR = NewPahoLogAdapter(zerolog.ErrorLevel)

	// broker connection string
	brokerURI, err := url.Parse(common.Config.MQTTBrokerUrl)
	if err != nil {
		return pkgerr.Wrap(err, "url.Parse broker url")
	}
	var pwd string
	if dc, ok := l.cache.dynamicConfig.ByKey(common.DynamicConfigMQTTAuth); ok {
		pwd = dc.Value
	}
	if pwd != "" {
		if brokerURI.User != nil {
			brokerURI.User = url.UserPassword(brokerURI.User.Username(), pwd)
		} else {
			brokerURI.User = url.UserPassword("gxydb-api", pwd)
		}
	}

	// client
	opts := mqtt.NewClientOptions().
		AddBroker(brokerURI.String()).
		SetClientID(common.Config.MQTTClientID).
		SetAutoReconnect(true).
		SetOnConnectHandler(l.Subscribe)
	l.client = mqtt.NewClient(opts)

	// connect
	if token := l.client.Connect(); token.Wait() && token.Error() != nil {
		return pkgerr.Wrap(token.Error(), "mqtt.client Connect")
	}

	// Start periodic admin messages
	go l.startPeriodicAdminMessages()

	return nil
}

func (l *MQTTListener) Subscribe(c mqtt.Client) {
	if token := l.client.Subscribe("galaxy/service/#", byte(2), l.HandleServiceProtocol); token.Wait() && token.Error() != nil {
		log.Error().Err(token.Error()).Msg("mqtt.client Subscribe galaxy/service")
	}
	// We use mqtt broker filter to pass only needed events, so we use qos 1 here
	if token := l.client.Subscribe("gxydb/events/#", byte(1), l.HandleEvent); token.Wait() && token.Error() != nil {
		log.Error().Err(token.Error()).Msg("mqtt.client Subscribe gxydb/events")
	}
	if token := l.client.Subscribe("gxydb/users/#", byte(1), l.UpdateSession); token.Wait() && token.Error() != nil {
		log.Error().Err(token.Error()).Msg("mqtt.client Subscribe gxydb/users")
	}

	// Subscribe to Janus gateway status messages
	statusTopic := "janus/+/status"
	if token := l.client.Subscribe(statusTopic, byte(1), l.HandleGatewayStatus); token.Wait() && token.Error() != nil {
		log.Error().Err(token.Error()).Str("topic", statusTopic).Msg("mqtt.client Subscribe status")
	} else {
		log.Info().Str("topic", statusTopic).Msg("Subscribed to gateway status")
	}

	// Subscribe to Janus admin responses
	adminTopic := "janus/+/from-janus-admin"
	if token := l.client.Subscribe(adminTopic, byte(1), l.HandleGatewayAdminResponse); token.Wait() && token.Error() != nil {
		log.Error().Err(token.Error()).Str("topic", adminTopic).Msg("mqtt.client Subscribe admin")
	} else {
		log.Info().Str("topic", adminTopic).Msg("Subscribed to gateway admin responses")
	}
}

func (l *MQTTListener) Close() {
	if l.periodicTicker != nil {
		l.periodicTicker.Stop()
	}
	l.client.Disconnect(1000)
}

func (l *MQTTListener) HandleServiceProtocol(c mqtt.Client, m mqtt.Message) {
	log.Debug().
		Bool("Duplicate", m.Duplicate()).
		Int8("QOS", int8(m.Qos())).
		Bool("Retained", m.Retained()).
		Str("Topic", m.Topic()).
		Uint16("MessageID", m.MessageID()).
		Bytes("payload", m.Payload()).
		Msg("MQTT handle service protocol")

	// A MessageHandler (called when a new message is received) must not block (unless ClientOptions.SetOrderMatters(false) set). If you wish to perform a long-running task, or publish a message, then please use a go routine (blocking in the handler is a common cause of unexpected pingresp  not received, disconnecting errors).
	go func() {
		if err := l.serviceProtocolHandler.HandleMessage(string(m.Payload())); err != nil {
			log.Error().Err(err).Msg("service protocol error")
		}
	}()
}

func (l *MQTTListener) HandleEvent(c mqtt.Client, m mqtt.Message) {
	log.Debug().
		Bool("Duplicate", m.Duplicate()).
		Int8("QOS", int8(m.Qos())).
		Bool("Retained", m.Retained()).
		Str("Topic", m.Topic()).
		Uint16("MessageID", m.MessageID()).
		Bytes("payload", m.Payload()).
		Msg("MQTT handle event")

	ctx := context.Background()
	event, err := janus.ParseEvent(m.Payload())
	if err != nil {
		log.Error().Err(err).Msg("parsing event error")
		return
	}

	go func() {
		if err := l.SessionManager.HandleEvent(ctx, event); err != nil {
			log.Error().Err(err).Msg("event error")
		}
	}()
}

func (l *MQTTListener) UpdateSession(c mqtt.Client, m mqtt.Message) {
	log.Debug().
		Bool("Duplicate", m.Duplicate()).
		Int8("QOS", int8(m.Qos())).
		Bool("Retained", m.Retained()).
		Str("Topic", m.Topic()).
		Uint16("MessageID", m.MessageID()).
		Bytes("payload", m.Payload()).
		Msg("MQTT update user session")
	var user *V1User
	if err := json.Unmarshal(m.Payload(), &user); err != nil {
		log.Error().Err(err).Msg("json.Unmarshal")
		return
	}
	ctx := context.Background()

	go func() {
		if err := l.SessionManager.UpsertSession(ctx, user); err != nil {
			log.Error().Err(err).Msg("update session error")
		}
	}()
}

// startPeriodicAdminMessages sends list_sessions requests to online gateways every 10 seconds
func (l *MQTTListener) startPeriodicAdminMessages() {
	l.periodicTicker = time.NewTicker(10 * time.Second)
	defer l.periodicTicker.Stop()

	for range l.periodicTicker.C {
		l.gatewayStatusesMu.RLock()
		for _, status := range l.gatewayStatuses {
			status.mu.RLock()
			online := status.Online
			name := status.Name
			status.mu.RUnlock()

			if online {
				topic := fmt.Sprintf("janus/%s/to-janus-admin", name)
				l.SendAdminMessage(topic)
			}
		}
		l.gatewayStatusesMu.RUnlock()
	}
}

// SendAdminMessage sends a list_sessions request to a Janus gateway
func (l *MQTTListener) SendAdminMessage(topic string) {
	message := map[string]interface{}{
		"janus":        "list_sessions",
		"transaction":  "transaction",
		"admin_secret": l.adminSecret,
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Error().Err(err).Msg("SendAdminMessage: failed to marshal message")
		return
	}

	log.Debug().Str("topic", topic).Bytes("message", jsonMessage).Msg("SendAdminMessage")

	if token := l.client.Publish(topic, byte(1), false, jsonMessage); token.Wait() && token.Error() != nil {
		log.Error().Err(token.Error()).Str("topic", topic).Msg("SendAdminMessage: failed to publish")
	}
}

// JanusStatusMessage represents the status message from Janus gateway
type JanusStatusMessage struct {
	Online bool `json:"online"`
}

// HandleGatewayStatus processes status messages from Janus gateways (janus/{server}/status)
func (l *MQTTListener) HandleGatewayStatus(c mqtt.Client, m mqtt.Message) {
	go func() {
		// Extract server name from topic: janus/{server}/status
		parts := strings.Split(m.Topic(), "/")
		if len(parts) < 2 {
			log.Error().Str("topic", m.Topic()).Msg("HandleGatewayStatus: invalid topic format")
			return
		}

		// Check if this is a Janus gateway (gxy*)
		serverName := parts[1]
		matched, _ := regexp.MatchString(`gxy\d+`, serverName)
		if !matched {
			log.Debug().Str("server", serverName).Msg("HandleGatewayStatus: ignoring non-gateway server")
			return
		}

		log.Debug().Str("topic", m.Topic()).Bytes("payload", m.Payload()).Msg("HandleGatewayStatus")

		var statusMsg JanusStatusMessage
		if err := json.Unmarshal(m.Payload(), &statusMsg); err != nil {
			log.Error().Err(err).Str("topic", m.Topic()).Msg("HandleGatewayStatus: failed to unmarshal")
			return
		}

		// Update or create gateway status
		l.gatewayStatusesMu.Lock()
		status, exists := l.gatewayStatuses[serverName]
		if !exists {
			status = &GatewayStatus{
				Name: serverName,
			}
			l.gatewayStatuses[serverName] = status
		}
		l.gatewayStatusesMu.Unlock()

		status.mu.Lock()
		status.Online = statusMsg.Online
		status.LastSeen = time.Now()
		status.mu.Unlock()

		log.Info().
			Str("server", serverName).
			Bool("online", statusMsg.Online).
			Msg("Gateway status updated")
	}()
}

// JanusAdminResponse represents the response from Janus admin API
type JanusAdminResponse struct {
	Janus       string  `json:"janus"`
	Transaction string  `json:"transaction"`
	Sessions    []int64 `json:"sessions"`
}

// HandleGatewayAdminResponse processes admin responses from Janus gateways (janus/{server}/from-janus-admin)
func (l *MQTTListener) HandleGatewayAdminResponse(c mqtt.Client, m mqtt.Message) {
	go func() {
		// Extract server name from topic: janus/{server}/from-janus-admin
		parts := strings.Split(m.Topic(), "/")
		if len(parts) < 2 {
			log.Error().Str("topic", m.Topic()).Msg("HandleGatewayAdminResponse: invalid topic format")
			return
		}

		serverName := parts[1]
		log.Debug().Str("topic", m.Topic()).Bytes("payload", m.Payload()).Msg("HandleGatewayAdminResponse")

		var response JanusAdminResponse
		if err := json.Unmarshal(m.Payload(), &response); err != nil {
			log.Error().Err(err).Str("topic", m.Topic()).Msg("HandleGatewayAdminResponse: failed to unmarshal")
			return
		}

		if response.Janus == "success" {
			// Update sessions count
			l.gatewayStatusesMu.RLock()
			status, exists := l.gatewayStatuses[serverName]
			l.gatewayStatusesMu.RUnlock()

			if exists {
				status.mu.Lock()
				status.Sessions = len(response.Sessions)
				status.mu.Unlock()

				log.Debug().
					Str("server", serverName).
					Int("sessions", len(response.Sessions)).
					Msg("Gateway sessions updated")
			}
		}
	}()
}

// GetGatewayStatuses returns a copy of all gateway statuses (thread-safe)
// Returns simple structs without mutexes for external use
func (l *MQTTListener) GetGatewayStatuses() map[string]*common.GatewayStatusInfo {
	l.gatewayStatusesMu.RLock()
	defer l.gatewayStatusesMu.RUnlock()

	result := make(map[string]*common.GatewayStatusInfo)
	for name, status := range l.gatewayStatuses {
		status.mu.RLock()
		result[name] = &common.GatewayStatusInfo{
			Name:     status.Name,
			Online:   status.Online,
			Sessions: status.Sessions,
			LastSeen: status.LastSeen,
		}
		status.mu.RUnlock()
	}

	return result
}

type PahoLogAdapter struct {
	level zerolog.Level
}

func NewPahoLogAdapter(level zerolog.Level) *PahoLogAdapter {
	return &PahoLogAdapter{level: level}
}

func (a *PahoLogAdapter) Println(v ...interface{}) {
	log.WithLevel(a.level).Msgf("mqtt: %s", fmt.Sprint(v...))
}

func (a *PahoLogAdapter) Printf(format string, v ...interface{}) {
	log.WithLevel(a.level).Msgf("mqtt: %s", fmt.Sprintf(format, v...))
}
