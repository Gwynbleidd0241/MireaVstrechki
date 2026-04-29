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
	agendaService *service.AgendaService,
) *Server {
	mux := http.NewServeMux()

	userHandler := handlers.NewUserHandler(userService, eventService, taskService, logger)
	eventHandler := handlers.NewEventHandler(eventService, logger)
	taskHandler := handlers.NewTaskHandler(taskService, logger)
	participantHandler := handlers.NewParticipantHandler(participantService, logger)
	agendaHandler := handlers.NewAgendaHandler(agendaService, logger)

	mux.Handle("/register", http.HandlerFunc(userHandler.Register))
	mux.Handle("/login", http.HandlerFunc(userHandler.Login))
	mux.Handle("/me", middleware.AuthMiddleware(cfg.Auth.JWTSecret, http.HandlerFunc(userHandler.Me)))
	mux.Handle("/me/events", middleware.AuthMiddleware(cfg.Auth.JWTSecret, http.HandlerFunc(userHandler.MyEvents)))
	mux.Handle("/me/tasks", middleware.AuthMiddleware(cfg.Auth.JWTSecret, http.HandlerFunc(userHandler.MyTasks)))
	createEventHandler := middleware.RequireRole("admin", "organizer")(
		http.HandlerFunc(eventHandler.Create),
	)
	mux.Handle("/events", middleware.AuthMiddleware(cfg.Auth.JWTSecret, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			eventHandler.List(w, r)
		case http.MethodPost:
			createEventHandler.ServeHTTP(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})),
	)
	mux.Handle("/users", middleware.AuthMiddleware(cfg.Auth.JWTSecret, http.HandlerFunc(userHandler.List)))
	mux.Handle("/events/", middleware.AuthMiddleware(cfg.Auth.JWTSecret, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/tasks/"):
			switch r.Method {
			case http.MethodGet:
				taskHandler.Get(w, r)
			case http.MethodPatch:
				taskHandler.Update(w, r)
			case http.MethodDelete:
				taskHandler.Delete(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}

		case strings.HasSuffix(r.URL.Path, "/tasks"):
			switch r.Method {
			case http.MethodGet:
				taskHandler.List(w, r)
			case http.MethodPost:
				taskHandler.Create(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}

		case strings.Contains(r.URL.Path, "/participants/"):
			switch r.Method {
			case http.MethodPatch:
				participantHandler.Update(w, r)
			case http.MethodDelete:
				participantHandler.Remove(w, r)
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

		case strings.Contains(r.URL.Path, "/agenda/"):
			switch r.Method {
			case http.MethodPatch:
				agendaHandler.Update(w, r)
			case http.MethodDelete:
				agendaHandler.Delete(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}

		case strings.HasSuffix(r.URL.Path, "/agenda"):
			switch r.Method {
			case http.MethodGet:
				agendaHandler.List(w, r)
			case http.MethodPost:
				agendaHandler.Add(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}

		default:
			switch r.Method {
			case http.MethodGet:
				eventHandler.Get(w, r)
			case http.MethodPatch:
				eventHandler.Update(w, r)
			case http.MethodDelete:
				eventHandler.Delete(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
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
