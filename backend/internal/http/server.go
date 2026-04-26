package http

import (
	"log"
	"net/http"

	"meeting-service/backend/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func New(cfg *config.Config) *Server {
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: mux,
	}

	return &Server{
		httpServer: srv,
	}
}

func (s *Server) Run() {
	log.Printf("server started on %s", s.httpServer.Addr)

	if err := s.httpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
