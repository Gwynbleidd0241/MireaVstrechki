package main

import (
	"log"
	"os"

	"github.com/pressly/goose/v3"

	"meeting-service/internal/config"
	httpserver "meeting-service/internal/http"
	"meeting-service/internal/logger"
	"meeting-service/internal/repository/postgres"
	"meeting-service/internal/service"
	"meeting-service/migrations"
)

func main() {
	cfg := config.Load()

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		runMigrations(cfg.Postgres.DSN)
		return
	}

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
	taskService := service.NewTaskService(taskRepo, eventRepo)

	participantRepo := postgres.NewParticipantRepository(db)
	participantService := service.NewParticipantService(participantRepo, eventRepo)

	agendaRepo := postgres.NewAgendaRepository(db)
	agendaService := service.NewAgendaService(agendaRepo)

	server := httpserver.New(cfg, logg, userService, eventService, taskService, participantService, agendaService)
	server.Run()
}

func runMigrations(dsn string) {
	db, err := postgres.New(dsn)
	if err != nil {
		log.Fatalf("migrate: connect: %v", err)
	}
	defer db.Close()

	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("migrate: dialect: %v", err)
	}

	if err := goose.Up(db, "."); err != nil {
		log.Fatalf("migrate: up: %v", err)
	}

	log.Println("migrations applied")
}
