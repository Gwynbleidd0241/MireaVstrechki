package config

import "os"

type Config struct {
	Server   Server
	Postgres Postgres
}

type Server struct {
	Port string
}

type Postgres struct {
	DSN string
}

func Load() *Config {
	cfg := defaultConfig()

	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.Server.Port = port
	}

	if dsn := os.Getenv("POSTGRES_DSN"); dsn != "" {
		cfg.Postgres.DSN = dsn
	}

	return cfg
}

func defaultConfig() *Config {
	return &Config{
		Server: Server{
			Port: "8080",
		},
		Postgres: Postgres{
			DSN: "postgres://postgres:postgres@localhost:5432/meeting_service?sslmode=disable",
		},
	}
}
