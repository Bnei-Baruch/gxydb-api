package instrumentation

import (
	"time"

	pkgerr "github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/queries"

	"github.com/Bnei-Baruch/gxydb-api/common"
)

// GatewayStatusProvider is an interface for getting gateway statuses
type GatewayStatusProvider interface {
	GetGatewayStatuses() map[string]*common.GatewayStatusInfo
}

type PeriodicCollector struct {
	ticker                *time.Ticker
	ticks                 int64
	db                    common.DBInterface
	gatewayStatusProvider GatewayStatusProvider
}

func NewPeriodicCollector(db common.DBInterface, gatewayStatusProvider GatewayStatusProvider) *PeriodicCollector {
	return &PeriodicCollector{
		ticker:                time.NewTicker(time.Second),
		db:                    db,
		gatewayStatusProvider: gatewayStatusProvider,
	}
}

func (pc *PeriodicCollector) Start() {
	if pc.ticker != nil {
		pc.ticker.Stop()
	}

	log.Info().Msg("periodically collecting stats")
	pc.ticker = time.NewTicker(time.Second)
	go pc.run()
}

func (pc *PeriodicCollector) Close() {
	if pc.ticker != nil {
		pc.ticker.Stop()
	}
}

func (pc *PeriodicCollector) run() {
	for range pc.ticker.C {
		pc.ticks++
		pc.collectRoomParticipants()
		pc.collectGatewaySessions()
	}
}

func (pc *PeriodicCollector) collectRoomParticipants() {
	rows, err := queries.Raw(`select r.name, count(distinct s.user_id)
										from sessions s inner join rooms r on s.room_id = r.gateway_uid
										where s.removed_at is null
										group by r.id;`).Query(pc.db)
	if err != nil {
		log.Error().Err(err).Msg("PeriodicCollector.collectRoomParticipants queries.Raw")
		return
	}

	Stats.RoomParticipantsGauge.Reset()

	for rows.Next() {
		var name string
		var count int64
		if err = rows.Scan(&name, &count); err != nil {
			log.Error().Err(err).Msg("PeriodicCollector.collectRoomParticipants rows.Scan")
		} else {
			Stats.RoomParticipantsGauge.WithLabelValues(name).Set(float64(count))
		}
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("PeriodicCollector.collectRoomParticipants rows.Err")
		return
	}
}

func (pc *PeriodicCollector) collectGatewaySessions() {
	Stats.GatewaySessionsGauge.Reset()

	// If MQTT is enabled, use data from MQTT
	if pc.gatewayStatusProvider != nil {
		statuses := pc.gatewayStatusProvider.GetGatewayStatuses()

		// Monitor AVAILABLE and FAILOVER servers (rooms type)
		allRoomsServers := make([]string, 0, len(common.Config.AvailableJanusServers)+len(common.Config.FailoverJanusServers))
		allRoomsServers = append(allRoomsServers, common.Config.AvailableJanusServers...)
		allRoomsServers = append(allRoomsServers, common.Config.FailoverJanusServers...)
		
		for _, serverName := range allRoomsServers {
			if status, ok := statuses[serverName]; ok {
				// Use data from MQTT - rooms type
				Stats.GatewaySessionsGauge.WithLabelValues(serverName, common.GatewayTypeRooms).Set(float64(status.Sessions))

				// Log if gateway is offline or stale
				if !status.Online {
					log.Debug().Str("gateway", serverName).Msg("Gateway is offline")
				} else if time.Since(status.LastSeen) > 60*time.Second {
					// Only warn if REALLY stale (>1 minute)
					log.Warn().
						Str("gateway", serverName).
						Dur("since_last_seen", time.Since(status.LastSeen)).
						Msg("Gateway status is stale")
				}
			} else {
				// Gateway not found in MQTT statuses - set to 0
				Stats.GatewaySessionsGauge.WithLabelValues(serverName, common.GatewayTypeRooms).Set(0)
				log.Debug().
					Str("gateway", serverName).
					Msg("Gateway not found in MQTT statuses")
			}
		}
		
		// Monitor streaming servers (streaming type)
		// FIXME: This is temporary monitoring for streaming servers.
		//        These servers should be monitored by strdb service instead.
		//        Move this monitoring logic to strdb and remove from gxydb-api.
		for _, serverName := range common.Config.StrJanusServers {
			if status, ok := statuses[serverName]; ok {
				// Use data from MQTT - streaming type
				Stats.GatewaySessionsGauge.WithLabelValues(serverName, common.GatewayTypeStreaming).Set(float64(status.Sessions))

				// Log if gateway is offline or stale
				if !status.Online {
					log.Debug().Str("gateway", serverName).Str("type", "streaming").Msg("Gateway is offline")
				} else if time.Since(status.LastSeen) > 60*time.Second {
					log.Warn().
						Str("gateway", serverName).
						Str("type", "streaming").
						Dur("since_last_seen", time.Since(status.LastSeen)).
						Msg("Gateway status is stale")
				}
			} else {
				// Gateway not found in MQTT statuses - set to 0
				Stats.GatewaySessionsGauge.WithLabelValues(serverName, common.GatewayTypeStreaming).Set(0)
				log.Debug().
					Str("gateway", serverName).
					Str("type", "streaming").
					Msg("Gateway not found in MQTT statuses")
			}
		}

		return
	}

	// MQTT is disabled - fall back to HTTP Admin API (legacy behavior)
	log.Warn().Msg("MQTT disabled - using legacy HTTP Admin API for gateway sessions")

	type gatewayCallRes struct {
		serverName string
		sessions   int
		duration   time.Duration
		err        error
	}

	c := make(chan *gatewayCallRes)

	// Use AVAILABLE + FAILOVER servers from config
	allServers := make([]string, 0, len(common.Config.AvailableJanusServers)+len(common.Config.FailoverJanusServers))
	allServers = append(allServers, common.Config.AvailableJanusServers...)
	allServers = append(allServers, common.Config.FailoverJanusServers...)
	
	for _, serverName := range allServers {
		go func(name string, c chan *gatewayCallRes) {
			res := &gatewayCallRes{serverName: name}
			start := time.Now()
			defer func() {
				res.duration = time.Since(start)
				c <- res
			}()

			log.Warn().
				Str("gateway", name).
				Msg("HTTP Admin API is deprecated - please enable MQTT")

			// Note: This requires janus_admin import which we removed
			// Leaving this as fallback but with error
			res.err = pkgerr.New("HTTP Admin API support removed - MQTT required")
		}(serverName, c)
	}

	timeout := time.After(900 * time.Millisecond)
	for i := range allServers {
		select {
		case res := <-c:
			if res.err != nil {
				log.Error().
					Err(res.err).
					Dur("duration", res.duration).
					Str("gateway", res.serverName).
					Msg("PeriodicCollector.collectGatewaySessions error")
			}
			Stats.GatewaySessionsGauge.WithLabelValues(res.serverName, common.GatewayTypeRooms).Set(float64(res.sessions))
		case <-timeout:
			log.Error().Msgf("PeriodicCollector.collectGatewaySessions timeout (i, len)=(%d,%d)", i, len(allServers))
			break
		}
	}
}
