package config

import (
	"os"
)

type Config struct {
	Server Server
}

type Server struct {
	Port string
}

func Load() *Config {
	cfg := defaultConfig()
	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.Server.Port = port
	}
	return cfg
}

func defaultConfig() *Config {
	return &Config{
		Server: Server{
			Port: "8080",
		},
	}
}
