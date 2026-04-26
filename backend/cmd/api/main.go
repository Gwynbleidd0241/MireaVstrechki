package main

import (
	"log"

	"meeting-service/internal/config"
	httpserver "meeting-service/internal/http"
	"meeting-service/internal/logger"
	"meeting-service/internal/repository/postgres"
)

func main() {
	cfg := config.Load()

	logg, err := logger.New()
	if err != nil {
		log.Fatal(err)
	}
	defer logg.Sync()

	db, err := postgres.New(cfg.Postgres.DSN)
	if err != nil {
		logg.Fatal("failed to connect to postgres")
	}
	defer db.Close()

	userRepo := postgres.NewUserRepository(db)

	server := httpserver.New(cfg, logg, userRepo)
	server.Run()
}
