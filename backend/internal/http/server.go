package http

import (
	"net/http"

	"go.uber.org/zap"

	"meeting-service/internal/config"
	"meeting-service/internal/http/handlers"
	"meeting-service/internal/http/middleware"
	"meeting-service/internal/service"
)

type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
}

func New(
	cfg *config.Config,
	logger *zap.Logger,
	userService *service.UserService,
	eventService *service.EventService,
) *Server {
	mux := http.NewServeMux()

	userHandler := handlers.NewUserHandler(userService, logger)
	eventHandler := handlers.NewEventHandler(eventService, logger)

	mux.HandleFunc("/register", userHandler.Register)
	mux.HandleFunc("/login", userHandler.Login)

	mux.Handle(
		"/me",
		middleware.AuthMiddleware(cfg.Auth.JWTSecret, http.HandlerFunc(userHandler.Me)),
	)

	mux.Handle(
		"/events",
		middleware.AuthMiddleware(cfg.Auth.JWTSecret, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				eventHandler.List(w, r)
			case http.MethodPost:
				eventHandler.Create(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		})),
	)

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
