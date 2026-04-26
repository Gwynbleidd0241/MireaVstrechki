package http

import (
	"net/http"
	"strings"

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
	taskService *service.TaskService,
	participantService *service.ParticipantService,
) *Server {
	mux := http.NewServeMux()

	userHandler := handlers.NewUserHandler(userService, logger)
	eventHandler := handlers.NewEventHandler(eventService, logger)
	taskHandler := handlers.NewTaskHandler(taskService, logger)
	participantHandler := handlers.NewParticipantHandler(participantService, logger)

	mux.Handle("/register", http.HandlerFunc(userHandler.Register))
	mux.Handle("/login", http.HandlerFunc(userHandler.Login))

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

	mux.Handle(
		"/events/",
		middleware.AuthMiddleware(cfg.Auth.JWTSecret, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case strings.HasSuffix(r.URL.Path, "/tasks"):
				switch r.Method {
				case http.MethodGet:
					taskHandler.List(w, r)
				case http.MethodPost:
					taskHandler.Create(w, r)
				default:
					http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				}

			case strings.HasSuffix(r.URL.Path, "/participants"):
				switch r.Method {
				case http.MethodGet:
					participantHandler.List(w, r)
				case http.MethodPost:
					participantHandler.Add(w, r)
				default:
					http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				}

			default:
				http.NotFound(w, r)
			}
		})),
	)

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + cfg.Server.Port,
			Handler: middleware.CORSMiddleware(mux),
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
