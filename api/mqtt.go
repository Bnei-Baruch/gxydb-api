package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	pkgerr "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Bnei-Baruch/gxydb-api/common"
)

type MQTTListener struct {
	client                 mqtt.Client
	cache                  *AppCache
	serviceProtocolHandler ServiceProtocolHandler
	SessionManager         SessionManager
}

func NewMQTTListener(cache *AppCache, sph ServiceProtocolHandler, sm SessionManager) *MQTTListener {
	return &MQTTListener{
		cache:                  cache,
		serviceProtocolHandler: sph,
		SessionManager:         sm,
	}
}

func (l *MQTTListener) Start() error {
	// logging
	mqtt.DEBUG = NewPahoLogAdapter(zerolog.InfoLevel)
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

	return nil
}

func (l *MQTTListener) Subscribe(c mqtt.Client) {
	if token := l.client.Subscribe("galaxy/service/#", byte(2), l.HandleServiceProtocol); token.Wait() && token.Error() != nil {
		log.Error().Err(token.Error()).Msg("mqtt.client Subscribe")
	}
	if token := l.client.Subscribe("gxydb/users/#", byte(1), l.UpdateSession); token.Wait() && token.Error() != nil {
		log.Error().Err(token.Error()).Msg("mqtt.client Subscribe")
	}
}

func (l *MQTTListener) Close() {
	l.client.Disconnect(1000)
}

func (l *MQTTListener) HandleServiceProtocol(c mqtt.Client, m mqtt.Message) {
	log.Info().
		Bool("Duplicate", m.Duplicate()).
		Int8("QOS", int8(m.Qos())).
		Bool("Retained", m.Retained()).
		Str("Topic", m.Topic()).
		Uint16("MessageID", m.MessageID()).
		Bytes("payload", m.Payload()).
		Msg("MQTT handle service protocol")
	if err := l.serviceProtocolHandler.HandleMessage(string(m.Payload())); err != nil {
		log.Error().Err(err).Msg("service protocol error")
	} else {
		m.Ack()
	}
}

func (l *MQTTListener) UpdateSession(c mqtt.Client, m mqtt.Message) {
	log.Info().
		Bool("Duplicate", m.Duplicate()).
		Int8("QOS", int8(m.Qos())).
		Bool("Retained", m.Retained()).
		Str("Topic", m.Topic()).
		Uint16("MessageID", m.MessageID()).
		Bytes("payload", m.Payload()).
		Msg("MQTT handle service protocol")
	var user *V1User
	if err := json.Unmarshal(m.Payload(), &user); err != nil {
		log.Error().Err(err).Msg("json.Unmarshal")
		return
	}
	ctx := context.Background()
	if err := l.SessionManager.UpsertSession(ctx, user); err != nil {
		log.Error().Err(err).Msg("update session error")
	} else {
		m.Ack()
	}
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
