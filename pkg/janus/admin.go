package janus

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
)

const defaultRequestTimeout = 10 * time.Second

// MQTTAdminClient implements Janus Admin API over MQTT.
// Requests are published to janus/{server}/to-janus-admin,
// responses are received from janus/{server}/from-janus-admin.
type MQTTAdminClient struct {
	client      mqtt.Client
	adminSecret string
	timeout     time.Duration

	mu       sync.Mutex
	pending  map[string]chan json.RawMessage
}

func NewMQTTAdminClient(client mqtt.Client, adminSecret string) *MQTTAdminClient {
	return &MQTTAdminClient{
		client:      client,
		adminSecret: adminSecret,
		timeout:     defaultRequestTimeout,
		pending:     make(map[string]chan json.RawMessage),
	}
}

// HandleResponse is the MQTT message handler for janus/+/from-janus-admin.
// It routes responses to pending requests by transaction ID.
func (c *MQTTAdminClient) HandleResponse(_ mqtt.Client, m mqtt.Message) {
	var base struct {
		Transaction string `json:"transaction"`
	}
	if err := json.Unmarshal(m.Payload(), &base); err != nil {
		log.Error().Err(err).Bytes("payload", m.Payload()).Msg("MQTTAdminClient: failed to unmarshal response")
		return
	}

	if base.Transaction == "" {
		return
	}

	c.mu.Lock()
	ch, ok := c.pending[base.Transaction]
	if ok {
		delete(c.pending, base.Transaction)
	}
	c.mu.Unlock()

	if ok {
		ch <- json.RawMessage(m.Payload())
	}
}

func (c *MQTTAdminClient) sendRequest(server string, payload map[string]interface{}) (json.RawMessage, error) {
	txID := generateTxID()
	payload["transaction"] = txID
	payload["admin_secret"] = c.adminSecret

	ch := make(chan json.RawMessage, 1)
	c.mu.Lock()
	c.pending[txID] = ch
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		delete(c.pending, txID)
		c.mu.Unlock()
	}()

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal request: %w", err)
	}

	topic := fmt.Sprintf("janus/%s/to-janus-admin", server)
	if token := c.client.Publish(topic, 1, false, b); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("mqtt publish to %s: %w", topic, token.Error())
	}

	select {
	case resp := <-ch:
		return resp, nil
	case <-time.After(c.timeout):
		return nil, fmt.Errorf("timeout waiting for response from %s (transaction %s)", server, txID)
	}
}

// MessagePlugin sends a message_plugin request to a Janus gateway.
// The request payload is sent as-is inside the "request" field.
func (c *MQTTAdminClient) MessagePlugin(server, plugin string, request map[string]interface{}) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"janus":   "message_plugin",
		"plugin":  plugin,
		"request": request,
	}

	raw, err := c.sendRequest(server, payload)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Janus    string                 `json:"janus"`
		Response map[string]interface{} `json:"response"`
		Error    *struct {
			Code   int    `json:"code"`
			Reason string `json:"reason"`
		} `json:"error"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("json.Unmarshal response: %w", err)
	}

	if resp.Janus == "error" && resp.Error != nil {
		return nil, fmt.Errorf("janus error [%d]: %s", resp.Error.Code, resp.Error.Reason)
	}

	if errReason, ok := resp.Response["error"]; ok {
		return nil, fmt.Errorf("plugin error: %v", errReason)
	}

	return resp.Response, nil
}

// HandleInfo sends a handle_info request to a Janus gateway.
func (c *MQTTAdminClient) HandleInfo(server string, sessionID, handleID uint64) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"janus":      "handle_info",
		"session_id": sessionID,
		"handle_id":  handleID,
	}

	raw, err := c.sendRequest(server, payload)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Janus string                 `json:"janus"`
		Info  map[string]interface{} `json:"info"`
		Error *struct {
			Code   int    `json:"code"`
			Reason string `json:"reason"`
		} `json:"error"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("json.Unmarshal response: %w", err)
	}

	if resp.Janus == "error" && resp.Error != nil {
		return nil, &AdminError{Code: resp.Error.Code, Reason: resp.Error.Reason}
	}

	return resp.Info, nil
}

// ListSessions sends a list_sessions request (fire-and-forget, no response expected).
func (c *MQTTAdminClient) ListSessions(server string) {
	payload := map[string]interface{}{
		"janus": "list_sessions",
	}
	payload["admin_secret"] = c.adminSecret
	payload["transaction"] = "transaction"

	b, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msg("MQTTAdminClient.ListSessions: marshal error")
		return
	}

	topic := fmt.Sprintf("janus/%s/to-janus-admin", server)
	if token := c.client.Publish(topic, 1, false, b); token.Wait() && token.Error() != nil {
		log.Error().Err(token.Error()).Str("topic", topic).Msg("MQTTAdminClient.ListSessions: publish error")
	}
}

// AdminError represents an error from Janus Admin API.
type AdminError struct {
	Code   int
	Reason string
}

func (e *AdminError) Error() string {
	return fmt.Sprintf("janus admin error [%d]: %s", e.Code, e.Reason)
}

func generateTxID() string {
	b := make([]byte, 12)
	rand.Read(b)
	return hex.EncodeToString(b)
}
