package main

import (
	"log"

	"meeting-service/internal/config"
	httpserver "meeting-service/internal/http"
	"meeting-service/internal/logger"
	"meeting-service/internal/repository/postgres"
	"meeting-service/internal/service"
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
	eventRepo := postgres.NewEventRepository(db)

	userService := service.NewUserService(userRepo, cfg.Auth.JWTSecret)
	eventService := service.NewEventService(eventRepo)

	taskRepo := postgres.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)

	participantRepo := postgres.NewParticipantRepository(db)
	participantService := service.NewParticipantService(participantRepo)

	server := httpserver.New(cfg, logg, userService, eventService, taskService, participantService)
	server.Run()
}
