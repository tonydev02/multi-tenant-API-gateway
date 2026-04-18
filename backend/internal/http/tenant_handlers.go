package gatewayhttp

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/namta/multi-tenant-api-gateway/backend/internal/auth"
	"github.com/namta/multi-tenant-api-gateway/backend/internal/tenant"
)

const timeLayout = "2006-01-02T15:04:05Z07:00"

type registerTenantRequest struct {
	TenantName string `json:"tenant_name"`
	TenantSlug string `json:"tenant_slug"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type updateTenantRequest struct {
	Name string `json:"name"`
}

func registerTenantHandler(tenantStore *tenant.Store, authStore *auth.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req registerTenantRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if req.TenantName == "" || req.TenantSlug == "" || req.Email == "" || req.Password == "" {
			writeError(w, http.StatusBadRequest, "tenant_name, tenant_slug, email, and password are required")
			return
		}

		t, err := tenantStore.Create(r.Context(), req.TenantName, req.TenantSlug)
		if err != nil {
			writeError(w, http.StatusBadRequest, "failed to create tenant")
			return
		}

		hash, err := auth.HashPassword(req.Password)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to hash password")
			return
		}
		admin, err := authStore.CreateAdminUser(r.Context(), t.ID, req.Email, hash)
		if err != nil {
			writeError(w, http.StatusBadRequest, "failed to create admin")
			return
		}

		writeJSON(w, http.StatusCreated, map[string]any{
			"tenant": t,
			"admin": map[string]any{
				"id":        admin.ID,
				"tenant_id": admin.TenantID,
				"email":     admin.Email,
			},
		})
	}
}

func getCurrentTenantHandler(tenantStore *tenant.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID, ok := auth.TenantIDFromContext(r.Context())
		if !ok {
			writeError(w, http.StatusUnauthorized, "tenant context missing")
			return
		}

		t, err := tenantStore.GetByID(r.Context(), tenantID)
		if shouldRetryReadError(err) {
			t, err = tenantStore.GetByID(r.Context(), tenantID)
		}
		if err != nil {
			requestID, _ := requestIDFromContext(r.Context())
			slog.Error("get current tenant failed",
				"request_id", requestID,
				"tenant_id", tenantID,
				"path", r.URL.Path,
				"error", err.Error(),
			)
			if errors.Is(err, sql.ErrNoRows) {
				writeError(w, http.StatusNotFound, "tenant not found")
				return
			}
			writeError(w, http.StatusInternalServerError, "failed to load tenant")
			return
		}
		writeJSON(w, http.StatusOK, t)
	}
}

func updateCurrentTenantHandler(tenantStore *tenant.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID, ok := auth.TenantIDFromContext(r.Context())
		if !ok {
			writeError(w, http.StatusUnauthorized, "tenant context missing")
			return
		}

		var req updateTenantRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if req.Name == "" {
			writeError(w, http.StatusBadRequest, "name is required")
			return
		}

		t, err := tenantStore.UpdateName(r.Context(), tenantID, req.Name)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeError(w, http.StatusNotFound, "tenant not found")
				return
			}
			writeError(w, http.StatusInternalServerError, "failed to update tenant")
			return
		}
		writeJSON(w, http.StatusOK, t)
	}
}

func deleteCurrentTenantHandler(tenantStore *tenant.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID, ok := auth.TenantIDFromContext(r.Context())
		if !ok {
			writeError(w, http.StatusUnauthorized, "tenant context missing")
			return
		}
		if err := tenantStore.Delete(r.Context(), tenantID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeError(w, http.StatusNotFound, "tenant not found")
				return
			}
			writeError(w, http.StatusInternalServerError, "failed to delete tenant")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	}
}

func consumerWhoAmIHandler(tenantStore *tenant.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID, ok := auth.TenantIDFromContext(r.Context())
		if !ok {
			writeError(w, http.StatusUnauthorized, "tenant context missing")
			return
		}
		t, err := tenantStore.GetByID(r.Context(), tenantID)
		if shouldRetryReadError(err) {
			t, err = tenantStore.GetByID(r.Context(), tenantID)
		}
		if err != nil {
			requestID, _ := requestIDFromContext(r.Context())
			slog.Error("consumer whoami tenant lookup failed",
				"request_id", requestID,
				"tenant_id", tenantID,
				"path", r.URL.Path,
				"error", err.Error(),
			)
			if errors.Is(err, sql.ErrNoRows) {
				writeError(w, http.StatusNotFound, "tenant not found")
				return
			}
			writeError(w, http.StatusInternalServerError, "failed to load tenant")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"tenant_id":   t.ID,
			"tenant_slug": t.Slug,
			"tenant_name": t.Name,
		})
	}
}
