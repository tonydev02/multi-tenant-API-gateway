package config

import (
	"fmt"
	"os"
	"strconv"
)

const defaultPort = 8080

// Config contains process configuration loaded from environment variables.
type Config struct {
	Port int
}

// Load reads config from environment and validates expected values.
func Load() (Config, error) {
	cfg := Config{Port: defaultPort}

	if rawPort, ok := os.LookupEnv("PORT"); ok && rawPort != "" {
		port, err := strconv.Atoi(rawPort)
		if err != nil {
			return Config{}, fmt.Errorf("parse PORT: %w", err)
		}
		if port <= 0 || port > 65535 {
			return Config{}, fmt.Errorf("PORT must be between 1 and 65535")
		}
		cfg.Port = port
	}

	return cfg, nil
}
