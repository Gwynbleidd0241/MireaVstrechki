// @title           Meeting Service API
// @version         1.0
// @description     Сервис для управления рабочими мероприятиями.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"

	"meeting-service/internal/config"
	httpserver "meeting-service/internal/http"
	"meeting-service/internal/logger"
	"meeting-service/internal/notification"
	"meeting-service/internal/repository/postgres"
	"meeting-service/internal/seed"
	"meeting-service/internal/service"
	"meeting-service/migrations"
)

func main() {
	cfg := config.Load()

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		runMigrations(cfg.Postgres.DSN)
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "seed" {
		runSeed(cfg.Postgres.DSN)
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
	agendaService := service.NewAgendaService(agendaRepo, eventRepo)

	reminderRepo := postgres.NewReminderRepository(db)
	emailSender := notification.NewEmailSender(cfg.SMTP)
	reminderScheduler := notification.NewReminderScheduler(reminderRepo, emailSender, logg)

	server := httpserver.New(cfg, logg, userService, eventService, taskService, participantService, agendaService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go server.Run()
	go reminderScheduler.Start(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logg.Fatal("server shutdown failed", zap.Error(err))
	}
	logg.Info("server stopped gracefully")
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

func runSeed(dsn string) {
	db, err := postgres.New(dsn)
	if err != nil {
		log.Fatalf("seed: connect: %v", err)
	}
	defer db.Close()

	if err := seed.Run(db); err != nil {
		log.Fatalf("seed: %v", err)
	}
}
