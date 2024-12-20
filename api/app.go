package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"

	"github.com/Bnei-Baruch/gxydb-api/common"
	"github.com/Bnei-Baruch/gxydb-api/domain"
	"github.com/Bnei-Baruch/gxydb-api/instrumentation"
	"github.com/Bnei-Baruch/gxydb-api/middleware"
)

type App struct {
	Router                 *mux.Router
	Handler                http.Handler
	DB                     common.DBInterface
	cache                  *AppCache
	sessionManager         SessionManager
	serviceProtocolHandler ServiceProtocolHandler
	gatewayTokensManager   *domain.GatewayTokensManager
	roomsStatisticsManager *domain.RoomStatisticsManager
	periodicStatsCollector *instrumentation.PeriodicCollector
	mqttListener           *MQTTListener
}

func (a *App) initOidc(issuerUrls []string) middleware.OIDCTokenVerifier {
	verifier, err := middleware.NewFailoverOIDCTokenVerifier(issuerUrls)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing OIDC token verifier")
	}

	return verifier
}

func (a *App) Initialize() {
	log.Info().Msg("initializing app")

	db, err := sql.Open("postgres", common.Config.DBUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("sql.Open")
	}
	db.SetMaxIdleConns(common.Config.DBMaxIdleConns)
	db.SetMaxOpenConns(common.Config.DBMaxOpenConns)
	db.SetConnMaxLifetime(common.Config.DBConnMaxLifetime)

	var tokenVerifier middleware.OIDCTokenVerifier
	if !common.Config.SkipAuth {
		tokenVerifier = a.initOidc(common.Config.AccountsUrls)
	}

	a.InitializeWithDeps(db, tokenVerifier)
}

func (a *App) InitializeWithDeps(db common.DBInterface, tokenVerifier middleware.OIDCTokenVerifier) {
	a.DB = db

	a.initRoutes()
	a.initCache()
	a.initSessionManagement()
	a.initGatewayTokensMonitoring()
	a.initRoomsStatistics()
	a.initServiceProtocolHandler()
	a.initMQTT()
	a.initInstrumentation()

	// this is declared here to abstract away the cache from auth middleware
	gatewayPwd := func(name string) (string, bool) {
		g, ok := a.cache.gateways.ByName(name)
		if ok {
			return g.EventsPassword, true
		}
		return "", false
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders: []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "Authorization"},
		MaxAge:         3600,
	})

	// middleware post route match
	a.Router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if route := mux.CurrentRoute(r); route != nil {
				rCtx, _ := middleware.ContextFromRequest(r)
				if tpl, err := route.GetPathTemplate(); err == nil {
					rCtx.RouteName = tpl
				} else {
					rCtx.RouteName = r.URL.Path
				}
			}
			h.ServeHTTP(w, r)
		})
	})

	a.Handler = middleware.ContextMiddleware(
		middleware.LoggingMiddleware(
			middleware.RecoveryMiddleware(
				middleware.RealIPMiddleware(
					corsMiddleware.Handler(
						middleware.AuthenticationMiddleware(tokenVerifier, gatewayPwd)(
							middleware.MinimalPermissionMiddleware(
								a.Router)))))))
}

func (a *App) Run() {
	defer a.Shutdown()

	addr := common.Config.ListenAddress
	log.Info().Msgf("app run %s", addr)
	if err := http.ListenAndServe(addr, a.Handler); err != nil {
		log.Fatal().Err(err).Msg("http.ListenAndServe")
	}
}

func (a *App) Shutdown() {
	if a.gatewayTokensManager != nil {
		a.gatewayTokensManager.Close()
	}
	if a.periodicStatsCollector != nil {
		a.periodicStatsCollector.Close()
	}
	if a.mqttListener != nil {
		a.mqttListener.Close()
	}
	a.sessionManager.Close()
	a.cache.Close()
	if err := a.DB.Close(); err != nil {
		log.Error().Err(err).Msg("DB.close")
	}
}

func (a *App) Notify(event interface{}) {
	switch event.(type) {
	case string:
		log.Info().Msgf("processing %s", event)
		switch event.(string) {
		case common.EventGatewayTokensChanged:
			if err := a.cache.gatewayTokens.Reload(a.DB); err != nil {
				log.Error().Err(err).Msg("cache.gatewayTokens.Reload")
			}
		}
	}
}

type promErrorLog struct {
	promhttp.Logger
}

func (l promErrorLog) Println(v ...interface{}) {
	log.Error().Msgf("prometheus metrics error: %s", fmt.Sprint(v...))
}

func (a *App) initRoutes() {
	a.Router = mux.NewRouter()

	// api v1 (current)
	a.Router.HandleFunc("/groups", a.V1ListGroups).Methods("GET")
	a.Router.HandleFunc("/group/{id}", a.V1CreateGroup).Methods("PUT") // unused
	a.Router.HandleFunc("/rooms", a.V1ListRooms).Methods("GET")
	a.Router.HandleFunc("/room/{id}", a.V1GetRoom).Methods("GET") // unused
	a.Router.HandleFunc("/rooms/{id}", a.V1UpdateRoom).Methods("PUT")
	a.Router.HandleFunc("/users", a.V1ListUsers).Methods("GET")    // unused
	a.Router.HandleFunc("/users/{id}", a.V1GetUser).Methods("GET") // unused
	a.Router.HandleFunc("/users/{id}", a.V1UpdateSession).Methods("PUT")
	a.Router.HandleFunc("/qids", a.V1ListComposites).Methods("GET")
	a.Router.HandleFunc("/qids/{id}", a.V1GetComposite).Methods("GET")
	a.Router.HandleFunc("/program/{id}", a.V1GetComposite).Methods("GET")
	a.Router.HandleFunc("/qids/{id}", a.V1UpdateComposite).Methods("PUT")

	a.Router.HandleFunc("/event", a.V1HandleEvent).Methods("POST")
	a.Router.HandleFunc("/protocol", a.V1HandleProtocol).Methods("POST")
	a.Router.HandleFunc("/protocol/service", a.V1HandleServiceProtocol).Methods("POST")

	// api v2 (next)
	a.Router.HandleFunc("/v2/config", a.V2GetConfig).Methods("GET")
	a.Router.HandleFunc("/v2/rooms_statistics", a.V2GetRoomsStatistics).Methods("GET") // Here due to more open permissions. otherwise might be under /admin/
	a.Router.HandleFunc("/v2/vhinfo", a.V2GetVHInfo).Methods("GET")

	// admin
	a.Router.HandleFunc("/admin/gateways", a.AdminListGateways).Methods("GET")
	a.Router.HandleFunc("/admin/gateways/{gateway_id}/sessions/{session_id}/handles/{handle_id}/info", a.AdminGatewaysHandleInfo).Methods("GET")
	a.Router.HandleFunc("/admin/rooms", a.AdminListRooms).Methods("GET")
	a.Router.HandleFunc("/admin/rooms", a.AdminCreateRoom).Methods("POST")
	a.Router.HandleFunc("/admin/rooms/{id}", a.AdminGetRoom).Methods("GET")
	a.Router.HandleFunc("/admin/rooms/{id}", a.AdminUpdateRoom).Methods("PUT")
	a.Router.HandleFunc("/admin/rooms/{id}", a.AdminDeleteRoom).Methods("DELETE")
	a.Router.HandleFunc("/admin/rooms_statistics", a.AdminDeleteRoomsStatistics).Methods("DELETE")
	a.Router.HandleFunc("/admin/dynamic_config", a.AdminListDynamicConfigs).Methods("GET")
	a.Router.HandleFunc("/admin/dynamic_config", a.AdminCreateDynamicConfig).Methods("POST")
	a.Router.HandleFunc("/admin/dynamic_config/{id}", a.AdminGetDynamicConfig).Methods("GET")
	a.Router.HandleFunc("/admin/dynamic_config/{id}", a.AdminUpdateDynamicConfig).Methods("PUT")
	a.Router.HandleFunc("/admin/dynamic_config/{key}", a.AdminSetDynamicConfig).Methods("POST")
	a.Router.HandleFunc("/admin/dynamic_config/{id}", a.AdminDeleteDynamicConfig).Methods("DELETE")

	// misc
	a.Router.HandleFunc("/health_check", a.HealthCheck).Methods("GET")
	a.Router.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			ErrorLog:          promErrorLog{},
			EnableOpenMetrics: true,
		},
	))
}

func (a *App) initCache() {
	a.cache = new(AppCache)
	if err := a.cache.Init(a.DB); err != nil {
		log.Fatal().Err(err).Msg("initialize app cache")
	}
}

func (a *App) initSessionManagement() {
	a.sessionManager = NewV1SessionManager(a.DB, a.cache)
	a.sessionManager.Start()
}

func (a *App) initServiceProtocolHandler() {
	a.serviceProtocolHandler = NewV1ServiceProtocolHandler(a.cache, a.roomsStatisticsManager)
}

func (a *App) initGatewayTokensMonitoring() {
	if common.Config.MonitorGatewayTokens {
		a.gatewayTokensManager = domain.NewGatewayTokensManager(a.DB, 3*24*time.Hour)
		a.gatewayTokensManager.AddObserver(a)
		a.gatewayTokensManager.Monitor()
	}
}

func (a *App) initRoomsStatistics() {
	a.roomsStatisticsManager = domain.NewRoomStatisticsManager(a.DB)
}

func (a *App) initInstrumentation() {
	instrumentation.Stats.Init()
	if common.Config.CollectPeriodicStats {
		a.periodicStatsCollector = instrumentation.NewPeriodicCollector(a.DB)
		a.periodicStatsCollector.Start()
	}
}

func (a *App) initMQTT() {
	if common.Config.MQTTBrokerUrl != "" {
		a.mqttListener = NewMQTTListener(a.cache, a.serviceProtocolHandler, a.sessionManager)
		if err := a.mqttListener.Start(); err != nil {
			log.Fatal().Err(err).Msg("initialize mqtt listener")
		}
	}
}
