package instrumentation

import (
	"github.com/prometheus/client_golang/prometheus"
)

var Stats = new(Collectors)

type Collectors struct {
	GatewaySessionsGauge     *prometheus.GaugeVec
	RoomParticipantsGauge    *prometheus.GaugeVec
	RequestDurationHistogram *prometheus.HistogramVec
	FailoverEventsCounter    *prometheus.CounterVec
	FailoverActiveGauge      *prometheus.GaugeVec
}

func (c *Collectors) Init() {
	c.GatewaySessionsGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "galaxy",
		Subsystem: "gateways",
		Name:      "sessions",
		Help:      "WebRTC Gateways active sessions",
	}, []string{
		// gateway name
		"name",
		// gateway type (rooms, streaming)
		"type"})

	c.RoomParticipantsGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "galaxy",
		Subsystem: "api",
		Name:      "participants",
		Help:      "Active room participants",
	}, []string{
		// room name
		"name",
	})

	c.RequestDurationHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "galaxy",
		Subsystem: "api",
		Name:      "request_duration",
		Help:      "Time (in milliseconds) spent serving HTTP requests.",
	}, []string{"method", "route", "status_code"})

	c.FailoverEventsCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "galaxy",
		Subsystem: "failover",
		Name:      "events_total",
		Help:      "Total number of failover events",
	}, []string{
		// event type: triggered, completed, failed, recovered
		"event",
		// failed server name
		"failed_server",
		// failover server name (empty for emergency distribution)
		"failover_server",
	})

	c.FailoverActiveGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "galaxy",
		Subsystem: "failover",
		Name:      "active",
		Help:      "Currently active failover mappings (1 = active, 0 = inactive)",
	}, []string{
		// failed server name
		"failed_server",
		// failover server name
		"failover_server",
	})

	prometheus.MustRegister(c.GatewaySessionsGauge)
	prometheus.MustRegister(c.RoomParticipantsGauge)
	prometheus.MustRegister(c.RequestDurationHistogram)
	prometheus.MustRegister(c.FailoverEventsCounter)
	prometheus.MustRegister(c.FailoverActiveGauge)
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())
}

func (c *Collectors) Reset() {
	c.GatewaySessionsGauge.Reset()
	c.RoomParticipantsGauge.Reset()
	c.RequestDurationHistogram.Reset()
	c.FailoverEventsCounter.Reset()
	c.FailoverActiveGauge.Reset()
}
