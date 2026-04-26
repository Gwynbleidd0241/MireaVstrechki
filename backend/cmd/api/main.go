package main

import (
	"log"
	"meeting-service/internal/config"
	httpserver "meeting-service/internal/http"
	"meeting-service/internal/logger"
)

func main() {
	cfg := config.Load()

	logg, err := logger.New()
	if err != nil {
		log.Fatal(err)
	}
	defer logg.Sync()

	server := httpserver.New(cfg, logg)
	server.Run()
}
