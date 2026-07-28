package config

import (
	"fmt"
	"strconv"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func Load(envGetter func(string) string) (*Config, error) {
	cfg := &Config{
		DatabaseURL: envGetter("DATABASE_URL"),
		Port:        orDefault(envGetter("PORT"), "8080"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("config: DATABASE_URL is required")
	}
	if _, err := strconv.Atoi(cfg.Port); err != nil {
		return nil, fmt.Errorf("config: PORT must be numeric, got %q", cfg.Port)
	}

	return cfg, nil
}

func orDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
