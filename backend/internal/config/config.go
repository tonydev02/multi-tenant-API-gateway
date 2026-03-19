package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultPort           = 8080
	defaultDBMaxOpenConns = 10
	defaultDBMaxIdleConns = 5
	defaultJWTExpiryMins  = 60
	minJWTSecretLen       = 32
)

// Config contains process configuration loaded from environment variables.
type Config struct {
	Environment     string
	Port            int
	FrontendOrigin  string
	DBURL           string
	DBMaxOpen       int
	DBMaxIdle       int
	RedisAddr       string
	RedisUsername   string
	RedisPassword   string
	RedisDB         int
	RedisTLS        bool
	RateLimitReqs   int64
	RateLimitWindow time.Duration
	ProxyTimeout    time.Duration
	ProxyUpstreams  string
	JWTSecret       string
	JWTIssuer       string
	JWTExpiry       time.Duration
	BootstrapOn     bool
	BootstrapName   string
	BootstrapSlug   string
	BootstrapMail   string
	BootstrapPass   string
}

// Load reads config from environment and validates expected values.
func Load() (Config, error) {
	cfg := Config{
		Environment:     getenv("ENVIRONMENT", "development"),
		Port:            defaultPort,
		FrontendOrigin:  getenv("FRONTEND_ORIGIN", "http://localhost:5173"),
		DBMaxOpen:       defaultDBMaxOpenConns,
		DBMaxIdle:       defaultDBMaxIdleConns,
		JWTIssuer:       getenv("JWT_ISSUER", "gateway-admin"),
		JWTExpiry:       time.Duration(defaultJWTExpiryMins) * time.Minute,
		BootstrapOn:     strings.EqualFold(getenv("BOOTSTRAP_ON_START", "true"), "true"),
		BootstrapName:   getenv("BOOTSTRAP_TENANT_NAME", "Acme"),
		BootstrapSlug:   getenv("BOOTSTRAP_TENANT_SLUG", "acme"),
		BootstrapMail:   getenv("BOOTSTRAP_ADMIN_EMAIL", "admin@acme.local"),
		BootstrapPass:   getenv("BOOTSTRAP_ADMIN_PASSWORD", "changeme123456"),
		RedisAddr:       getenv("REDIS_ADDR", fmt.Sprintf("%s:%s", getenv("REDIS_HOST", "127.0.0.1"), getenv("REDIS_PORT", "56379"))),
		RedisUsername:   getenv("REDIS_USERNAME", ""),
		RedisPassword:   getenv("REDIS_PASSWORD", ""),
		RedisTLS:        strings.EqualFold(getenv("REDIS_TLS", "false"), "true"),
		RateLimitReqs:   60,
		RateLimitWindow: 60 * time.Second,
		ProxyTimeout:    10 * time.Second,
		ProxyUpstreams:  getenv("PROXY_UPSTREAMS", ""),
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
	cfg.RedisDB, err = parseIntEnv("REDIS_DB", 0)
	if err != nil {
		return Config{}, err
	}
	rateLimitReqs, err := parseIntEnv("RATE_LIMIT_REQUESTS", int(cfg.RateLimitReqs))
	if err != nil {
		return Config{}, err
	}
	rateLimitWindowSeconds, err := parseIntEnv("RATE_LIMIT_WINDOW_SECONDS", int(cfg.RateLimitWindow.Seconds()))
	if err != nil {
		return Config{}, err
	}
	proxyTimeoutSeconds, err := parseIntEnv("PROXY_TIMEOUT_SECONDS", int(cfg.ProxyTimeout.Seconds()))
	if err != nil {
		return Config{}, err
	}
	cfg.RateLimitReqs = int64(rateLimitReqs)
	cfg.RateLimitWindow = time.Duration(rateLimitWindowSeconds) * time.Second
	cfg.ProxyTimeout = time.Duration(proxyTimeoutSeconds) * time.Second

	cfg.DBURL = getenv("DATABASE_URL", "postgres://gateway:gateway@localhost:5432/gateway?sslmode=disable")
	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}
	if len(cfg.JWTSecret) < minJWTSecretLen {
		return Config{}, fmt.Errorf("JWT_SECRET must be at least %d characters", minJWTSecretLen)
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
	if strings.TrimSpace(cfg.FrontendOrigin) == "" {
		return Config{}, fmt.Errorf("FRONTEND_ORIGIN is required")
	}
	if strings.TrimSpace(cfg.DBURL) == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.DBMaxOpen <= 0 || cfg.DBMaxIdle <= 0 {
		return Config{}, fmt.Errorf("DB_MAX_OPEN_CONNS and DB_MAX_IDLE_CONNS must be greater than zero")
	}
	if cfg.RateLimitReqs <= 0 || cfg.RateLimitWindow <= 0 {
		return Config{}, fmt.Errorf("RATE_LIMIT_REQUESTS and RATE_LIMIT_WINDOW_SECONDS must be greater than zero")
	}
	if cfg.ProxyTimeout <= 0 {
		return Config{}, fmt.Errorf("PROXY_TIMEOUT_SECONDS must be greater than zero")
	}
	if cfg.Environment != "development" && cfg.Environment != "staging" && cfg.Environment != "production" {
		return Config{}, fmt.Errorf("ENVIRONMENT must be one of development, staging, production")
	}
	if cfg.Environment != "development" && cfg.BootstrapOn {
		return Config{}, fmt.Errorf("BOOTSTRAP_ON_START must be false when ENVIRONMENT is %s", cfg.Environment)
	}
	if cfg.BootstrapOn {
		if strings.TrimSpace(cfg.BootstrapName) == "" || strings.TrimSpace(cfg.BootstrapSlug) == "" || strings.TrimSpace(cfg.BootstrapMail) == "" {
			return Config{}, fmt.Errorf("bootstrap tenant and admin identity fields are required when BOOTSTRAP_ON_START=true")
		}
		if len(cfg.BootstrapPass) < 12 {
			return Config{}, fmt.Errorf("BOOTSTRAP_ADMIN_PASSWORD must be at least 12 characters when BOOTSTRAP_ON_START=true")
		}
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
