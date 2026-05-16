package config

import "os"

type Config struct {
	Server   Server
	Postgres Postgres
	Auth     Auth
	SMTP     SMTP
	Yandex   Yandex
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

type SMTP struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type Yandex struct {
	GeoKey string
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

	if host := os.Getenv("SMTP_HOST"); host != "" {
		cfg.SMTP.Host = host
	}
	if port := os.Getenv("SMTP_PORT"); port != "" {
		cfg.SMTP.Port = port
	}
	if user := os.Getenv("SMTP_USERNAME"); user != "" {
		cfg.SMTP.Username = user
	}
	if pass := os.Getenv("SMTP_PASSWORD"); pass != "" {
		cfg.SMTP.Password = pass
	}
	if from := os.Getenv("SMTP_FROM"); from != "" {
		cfg.SMTP.From = from
	}

	if key := os.Getenv("YANDEX_GEO_KEY"); key != "" {
		cfg.Yandex.GeoKey = key
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
