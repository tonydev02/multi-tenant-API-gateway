package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/namta/multi-tenant-api-gateway/backend/internal/tenant"
)

// EnsureBootstrap inserts default tenant/admin if configured and missing.
func EnsureBootstrap(ctx context.Context, db *sql.DB, tenantName, tenantSlug, adminEmail, adminPassword string) error {
	tenantStore := tenant.NewStore(db)
	authStore := NewStore(db)

	t, err := tenantStore.GetBySlug(ctx, tenantSlug)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("lookup bootstrap tenant: %w", err)
		}
		t, err = tenantStore.Create(ctx, tenantName, tenantSlug)
		if err != nil {
			return fmt.Errorf("create bootstrap tenant: %w", err)
		}
	}

	_, err = authStore.GetAdminByEmail(ctx, adminEmail)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("lookup bootstrap admin: %w", err)
	}

	hash, err := HashPassword(adminPassword)
	if err != nil {
		return fmt.Errorf("hash bootstrap password: %w", err)
	}
	if _, err := authStore.CreateAdminUser(ctx, t.ID, adminEmail, hash); err != nil {
		return fmt.Errorf("create bootstrap admin: %w", err)
	}
	return nil
}
