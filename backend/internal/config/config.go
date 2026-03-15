package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	defaultPort           = 8080
	defaultDBMaxOpenConns = 10
	defaultDBMaxIdleConns = 5
	defaultJWTExpiryMins  = 60
)

// Config contains process configuration loaded from environment variables.
type Config struct {
	Port           int
	FrontendOrigin string
	DBURL          string
	DBMaxOpen      int
	DBMaxIdle      int
	JWTSecret      string
	JWTIssuer      string
	JWTExpiry      time.Duration
	BootstrapOn    bool
	BootstrapName  string
	BootstrapSlug  string
	BootstrapMail  string
	BootstrapPass  string
}

// Load reads config from environment and validates expected values.
func Load() (Config, error) {
	cfg := Config{
		Port:           defaultPort,
		FrontendOrigin: getenv("FRONTEND_ORIGIN", "http://localhost:5173"),
		DBMaxOpen:      defaultDBMaxOpenConns,
		DBMaxIdle:      defaultDBMaxIdleConns,
		JWTIssuer:      getenv("JWT_ISSUER", "gateway-admin"),
		JWTExpiry:      time.Duration(defaultJWTExpiryMins) * time.Minute,
		BootstrapOn:    getenv("BOOTSTRAP_ON_START", "true") == "true",
		BootstrapName:  getenv("BOOTSTRAP_TENANT_NAME", "Acme"),
		BootstrapSlug:  getenv("BOOTSTRAP_TENANT_SLUG", "acme"),
		BootstrapMail:  getenv("BOOTSTRAP_ADMIN_EMAIL", "admin@acme.local"),
		BootstrapPass:  getenv("BOOTSTRAP_ADMIN_PASSWORD", "changeme123"),
	}

	var err error
	cfg.Port, err = parseIntEnv("PORT", cfg.Port)
	if err != nil {
		return Config{}, err
	}
	cfg.DBMaxOpen, err = parseIntEnv("DB_MAX_OPEN_CONNS", cfg.DBMaxOpen)
	if err != nil {
		return Config{}, err
	}
	cfg.DBMaxIdle, err = parseIntEnv("DB_MAX_IDLE_CONNS", cfg.DBMaxIdle)
	if err != nil {
		return Config{}, err
	}

	cfg.DBURL = getenv("DATABASE_URL", "postgres://gateway:gateway@localhost:5432/gateway?sslmode=disable")
	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}

	if rawMins := os.Getenv("JWT_EXPIRY_MINUTES"); rawMins != "" {
		mins, convErr := strconv.Atoi(rawMins)
		if convErr != nil {
			return Config{}, fmt.Errorf("parse JWT_EXPIRY_MINUTES: %w", convErr)
		}
		if mins <= 0 {
			return Config{}, fmt.Errorf("JWT_EXPIRY_MINUTES must be greater than zero")
		}
		cfg.JWTExpiry = time.Duration(mins) * time.Minute
	}

	if cfg.Port <= 0 || cfg.Port > 65535 {
		return Config{}, fmt.Errorf("PORT must be between 1 and 65535")
	}
	if cfg.DBMaxOpen <= 0 || cfg.DBMaxIdle <= 0 {
		return Config{}, fmt.Errorf("DB_MAX_OPEN_CONNS and DB_MAX_IDLE_CONNS must be greater than zero")
	}

	return cfg, nil
}

func parseIntEnv(name string, fallback int) (int, error) {
	raw := os.Getenv(name)
	if raw == "" {
		return fallback, nil
	}
	val, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", name, err)
	}
	return val, nil
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
