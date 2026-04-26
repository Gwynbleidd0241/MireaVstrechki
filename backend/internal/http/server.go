package http

import (
	"net/http"

	"go.uber.org/zap"

	"meeting-service/internal/config"
	"meeting-service/internal/http/handlers"
	"meeting-service/internal/service"
)

type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
}

func New(cfg *config.Config, logger *zap.Logger, userService *service.UserService) *Server {
	mux := http.NewServeMux()

	userHandler := handlers.NewUserHandler(userService, logger)

	mux.HandleFunc("/register", userHandler.Register)
	mux.HandleFunc("/login", userHandler.Login)

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
