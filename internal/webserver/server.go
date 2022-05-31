package webserver

import (
	"context"
	"github.com/Mortimor1/mikromon-discovery/internal/config"
	"github.com/Mortimor1/mikromon-discovery/internal/webserver/handlers"
	"github.com/Mortimor1/mikromon-discovery/pkg/logging"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(cfg *config.Config) error {
	// Init logger
	logger := logging.GetLogger()

	// Init http router
	logger.Info("Create new router")
	router := mux.NewRouter()
	router.Use(handlers.Middleware)
	router.Use(handlers.LoggingMiddleware)

	reg := NewMetricsRegistry()
	router.Handle("/metrics", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))

	logger.Info("Register handlers")

	s.httpServer = &http.Server{
		Addr:           cfg.Http.BindIp + ":" + cfg.Http.Port,
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	// Start http server
	logger.Infof("Server listening on %s:%s", cfg.Http.BindIp, cfg.Http.Port)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
