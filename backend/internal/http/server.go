package http

import (
	"net/http"

	"go.uber.org/zap"

	"meeting-service/internal/config"
	"meeting-service/internal/http/handlers"
	"meeting-service/internal/repository/postgres"
)

type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
}

func New(cfg *config.Config, logger *zap.Logger, userRepo *postgres.UserRepository) *Server {
	mux := http.NewServeMux()

	registerHandler := handlers.NewRegisterHandler(userRepo, logger)

	mux.Handle("/register", registerHandler)

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + cfg.Server.Port,
			Handler: mux,
		},
		logger: logger,
	}
}

func (s *Server) Run() {
	s.logger.Info("server started", zap.String("addr", s.httpServer.Addr))

	if err := s.httpServer.ListenAndServe(); err != nil {
		s.logger.Fatal("server failed", zap.Error(err))
	}
}
