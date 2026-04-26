package http

import (
	"log"
	"net/http"

	"meeting-service/internal/config"

	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
}

func New(cfg *config.Config, logger *zap.Logger) *Server {
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: mux,
	}

	return &Server{
		httpServer: srv,
		logger:     logger,
	}
}

func (s *Server) Run() {
	s.logger.Info("server started", zap.String("addr", s.httpServer.Addr))
	if err := s.httpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
