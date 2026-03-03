package testutil

import "os"

type GatewayConfig struct {
	GatewayURL  string
	AdminURL    string
	AdminSecret string
}

type GatewayManager struct {
	Config *GatewayConfig
}

func (m *GatewayManager) Init() {
	m.Config = &GatewayConfig{
		GatewayURL:  "ws://localhost:8188/",
		AdminURL:    "http://localhost:7088/admin",
		AdminSecret: "janusoverlord",
	}

	if val := os.Getenv("TEST_GATEWAY_URL"); val != "" {
		m.Config.GatewayURL = val
	}
	if val := os.Getenv("TEST_GATEWAY_ADMIN_URL"); val != "" {
		m.Config.AdminURL = val
	}
	if val := os.Getenv("TEST_GATEWAY_ADMIN_SECRET"); val != "" {
		m.Config.AdminSecret = val
	}
}
