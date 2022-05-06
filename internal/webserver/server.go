package webserver

import (
	"context"
	"github.com/Mortimor1/mikromon-discovery/internal/config"
	"github.com/Mortimor1/mikromon-discovery/internal/subnet"
	"github.com/Mortimor1/mikromon-discovery/internal/webserver/handlers"
	"github.com/Mortimor1/mikromon-discovery/pkg/logging"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Server struct {
	httpServer       *http.Server
	subnetRepository *subnet.SubnetRepository
}

func (s *Server) Run(cfg *config.Config) error {
	// Init logger
	logger := logging.GetLogger()

	// Init http router
	logger.Info("Create new router")
	router := mux.NewRouter()
	router.Use(handlers.Middleware)
	router.Use(handlers.LoggingMiddleware)

	subnetHandler := subnet.NewSubnetHandler(logger, s.subnetRepository)

	logger.Info("Register handlers")
	subnetHandler.Register(router)

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

func NewHttpServer(subnetRepository *subnet.SubnetRepository) *Server {
	s := Server{
		httpServer:       nil,
		subnetRepository: subnetRepository,
	}
	return &s
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
