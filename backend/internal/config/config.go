package config

import "os"

type Config struct {
	Server   Server
	Postgres Postgres
	Auth     Auth
}

type Server struct {
	Port string
}

type Postgres struct {
	DSN string
}

type Auth struct {
	JWTSecret string
}

func Load() *Config {
	cfg := defaultConfig()

	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.Server.Port = port
	}

	if dsn := os.Getenv("POSTGRES_DSN"); dsn != "" {
		cfg.Postgres.DSN = dsn
	}

	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		cfg.Auth.JWTSecret = secret
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
		Auth: Auth{
			JWTSecret: "dev-secret",
		},
	}
}
