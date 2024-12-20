package common

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type config struct {
	ListenAddress         string
	DBUrl                 string
	AccountsUrls          []string
	SkipAuth              bool
	SkipEventsAuth        bool
	SkipPermissions       bool
	IceServers            map[string][]string
	ServicePasswords      []string
	Secret                string
	MonitorGatewayTokens  bool
	GatewayRoomsSecret    string
	GatewayPluginAdminKey string
	CollectPeriodicStats  bool
	CleanSessionsInterval time.Duration
	DeadSessionPeriod     time.Duration
	DBMaxIdleConns        int
	DBMaxOpenConns        int
	DBConnMaxLifetime     time.Duration
	MQTTBrokerUrl         string
	MQTTClientID          string
	MQTTPassword          string
	MQTTSecure            bool
	VHUrl                 string
}

func newConfig() *config {
	return &config{
		ListenAddress:         ":8081",
		DBUrl:                 "postgres://user:password@localhost/galaxy?sslmode=disable",
		AccountsUrls:          []string{"https://accounts.kab.info/auth/realms/main"},
		SkipAuth:              false,
		SkipEventsAuth:        false,
		SkipPermissions:       false,
		IceServers:            make(map[string][]string),
		ServicePasswords:      make([]string, 0),
		MonitorGatewayTokens:  true,
		GatewayRoomsSecret:    "",
		GatewayPluginAdminKey: "",
		CollectPeriodicStats:  true,
		CleanSessionsInterval: time.Minute,
		DeadSessionPeriod:     90 * time.Second,
		DBMaxIdleConns:        2,
		DBMaxOpenConns:        0,
		DBConnMaxLifetime:     0,
		MQTTBrokerUrl:         "",
		MQTTClientID:          "gxydb-api-dev",
		MQTTPassword:          "",
		MQTTSecure:            false,
		VHUrl:                 "https://api.kli.one",
	}
}

var Config *config

func Init() {
	Config = newConfig()

	if val := os.Getenv("LISTEN_ADDRESS"); val != "" {
		Config.ListenAddress = val
	}
	if val := os.Getenv("DB_URL"); val != "" {
		Config.DBUrl = val
	}
	if val := os.Getenv("ACCOUNTS_URL"); val != "" {
		Config.AccountsUrls = strings.Split(val, ",")
	}
	if val := os.Getenv("SKIP_AUTH"); val != "" {
		Config.SkipAuth = val == "true"
	}
	if val := os.Getenv("SKIP_EVENTS_AUTH"); val != "" {
		Config.SkipEventsAuth = val == "true"
	}
	if val := os.Getenv("SKIP_PERMISSIONS"); val != "" {
		Config.SkipPermissions = val == "true"
	}
	if val := os.Getenv("ICE_SERVERS_ROOMS"); val != "" {
		Config.IceServers["rooms"] = strings.Split(val, ",")
	}
	if val := os.Getenv("ICE_SERVERS_STREAMING"); val != "" {
		Config.IceServers["streaming"] = strings.Split(val, ",")
	}
	if val := os.Getenv("SERVICE_PASSWORDS"); val != "" {
		Config.ServicePasswords = strings.Split(val, ",")
	}
	if val := os.Getenv("SECRET"); val != "" {
		Config.Secret = val
	}
	if val := os.Getenv("MONITOR_GATEWAY_TOKENS"); val != "" {
		Config.MonitorGatewayTokens = val == "true"
	}
	if val := os.Getenv("GATEWAY_ROOMS_SECRET"); val != "" {
		Config.GatewayRoomsSecret = val
	}
	if val := os.Getenv("GATEWAY_PLUGIN_ADMIN_KEY"); val != "" {
		Config.GatewayPluginAdminKey = val
	}
	if val := os.Getenv("COLLECT_PERIODIC_STATS"); val != "" {
		Config.CollectPeriodicStats = val == "true"
	}
	if val := os.Getenv("CLEAN_SESSIONS_INTERVAL"); val != "" {
		pVal, err := time.ParseDuration(val)
		if err != nil {
			panic(err)
		}
		Config.CleanSessionsInterval = pVal
	}
	if val := os.Getenv("DEAD_SESSION_PERIOD"); val != "" {
		pVal, err := time.ParseDuration(val)
		if err != nil {
			panic(err)
		}
		if pVal <= 0 {
			panic(fmt.Errorf("DEAD_SESSION_PERIOD must be positive, got %d", pVal))
		}
		Config.DeadSessionPeriod = pVal
	}
	if val := os.Getenv("DB_MAX_IDLE_CONNS"); val != "" {
		pVal, err := strconv.Atoi(val)
		if err != nil {
			panic(err)
		}
		Config.DBMaxIdleConns = pVal
	}
	if val := os.Getenv("DB_MAX_OPEN_CONNS"); val != "" {
		pVal, err := strconv.Atoi(val)
		if err != nil {
			panic(err)
		}
		Config.DBMaxOpenConns = pVal
	}
	if val := os.Getenv("DB_CONN_MAX_LIFETIME"); val != "" {
		pVal, err := time.ParseDuration(val)
		if err != nil {
			panic(err)
		}
		Config.DBConnMaxLifetime = pVal
	}
	if val := os.Getenv("MQTT_BROKER_URL"); val != "" {
		Config.MQTTBrokerUrl = val
	}
	if val := os.Getenv("MQTT_CLIENT_ID"); val != "" {
		Config.MQTTClientID = val
	}
	if val := os.Getenv("MQTT_PASSWORD"); val != "" {
		Config.MQTTPassword = val
	}
	if val := os.Getenv("MQTT_SECURE"); val != "" {
		Config.MQTTSecure = val == "true"
	}
	if val := os.Getenv("VH_URL"); val != "" {
		Config.VHUrl = val
	}
}
