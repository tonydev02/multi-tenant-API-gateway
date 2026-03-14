package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/namta/multi-tenant-api-gateway/backend/internal/auth"
	"github.com/namta/multi-tenant-api-gateway/backend/internal/config"
	"github.com/namta/multi-tenant-api-gateway/backend/internal/db"
	gatewayhttp "github.com/namta/multi-tenant-api-gateway/backend/internal/http"
	"github.com/namta/multi-tenant-api-gateway/backend/internal/tenant"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	database, err := db.Open(context.Background(), cfg.DBURL, cfg.DBMaxOpen, cfg.DBMaxIdle)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer func() {
		if closeErr := database.Close(); closeErr != nil {
			log.Printf("close database: %v", closeErr)
		}
	}()

	if err := db.Migrate(context.Background(), database); err != nil {
		log.Fatalf("migrate database: %v", err)
	}
	if cfg.BootstrapOn {
		if err := auth.EnsureBootstrap(
			context.Background(),
			database,
			cfg.BootstrapName,
			cfg.BootstrapSlug,
			cfg.BootstrapMail,
			cfg.BootstrapPass,
		); err != nil {
			log.Fatalf("bootstrap auth data: %v", err)
		}
	}

	tenantStore := tenant.NewStore(database)
	authStore := auth.NewStore(database)
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTIssuer, cfg.JWTExpiry)

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),
		Handler: gatewayhttp.NewRouter(gatewayhttp.Dependencies{
			AuthStore:   authStore,
			TenantStore: tenantStore,
			JWTManager:  jwtManager,
			APIKeyAuth:  auth.NewAPIKeyAuthenticator(authStore),
		}),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if shutdownErr := server.Shutdown(shutdownCtx); shutdownErr != nil {
			log.Printf("server shutdown error: %v", shutdownErr)
		}
	}()

	log.Printf("backend listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
