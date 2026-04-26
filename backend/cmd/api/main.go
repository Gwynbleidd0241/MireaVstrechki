package main

import (
	"meeting-service/backend/internal/config"
	httpserver "meeting-service/backend/internal/http"
)

func main() {
	cfg := config.Load()

	server := httpserver.New(cfg)
	server.Run()
}
