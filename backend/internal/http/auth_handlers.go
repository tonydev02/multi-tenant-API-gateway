package gatewayhttp

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/namta/multi-tenant-api-gateway/backend/internal/auth"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type apiKeyCreateRequest struct {
	Name string `json:"name"`
}

type apiKeyCreateResponse struct {
	ID        int64  `json:"id"`
	TenantID  int64  `json:"tenant_id"`
	Name      string `json:"name"`
	Prefix    string `json:"prefix"`
	APIKey    string `json:"api_key"`
	CreatedAt string `json:"created_at"`
}

func loginHandler(store *auth.Store, jwt *auth.JWTManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if req.Email == "" || req.Password == "" {
			writeError(w, http.StatusBadRequest, "email and password are required")
			return
		}

		user, err := store.GetAdminByEmail(r.Context(), req.Email)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeError(w, http.StatusUnauthorized, "invalid credentials")
				return
			}
			writeError(w, http.StatusInternalServerError, "failed to load user")
			return
		}
		if err := auth.VerifyPassword(user.PasswordHash, req.Password); err != nil {
			writeError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}

		token, err := jwt.Issue(user)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to issue token")
			return
		}
		writeJSON(w, http.StatusOK, loginResponse{Token: token})
	}
}

func meHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := auth.ClaimsFromContext(r.Context())
		if !ok {
			writeError(w, http.StatusUnauthorized, "missing claims")
			return
		}
		writeJSON(w, http.StatusOK, claims)
	}
}

func createAPIKeyHandler(store *auth.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID, ok := auth.TenantIDFromContext(r.Context())
		if !ok {
			writeError(w, http.StatusUnauthorized, "tenant context missing")
			return
		}

		var req apiKeyCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if req.Name == "" {
			writeError(w, http.StatusBadRequest, "name is required")
			return
		}

		raw, prefix, hash, err := auth.GenerateAPIKey()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to generate api key")
			return
		}
		record, err := store.CreateAPIKey(r.Context(), tenantID, req.Name, prefix, hash)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to create api key")
			return
		}

		writeJSON(w, http.StatusCreated, apiKeyCreateResponse{
			ID:        record.ID,
			TenantID:  record.TenantID,
			Name:      record.Name,
			Prefix:    record.Prefix,
			APIKey:    raw,
			CreatedAt: record.CreatedAt.UTC().Format(timeLayout),
		})
	}
}

func listAPIKeysHandler(store *auth.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID, ok := auth.TenantIDFromContext(r.Context())
		if !ok {
			writeError(w, http.StatusUnauthorized, "tenant context missing")
			return
		}
		keys, err := store.ListAPIKeysByTenant(r.Context(), tenantID)
		if shouldRetryReadError(err) {
			keys, err = store.ListAPIKeysByTenant(r.Context(), tenantID)
		}
		if err != nil {
			requestID, _ := requestIDFromContext(r.Context())
			slog.Error("list api keys failed",
				"request_id", requestID,
				"tenant_id", tenantID,
				"path", r.URL.Path,
				"error", err.Error(),
			)
			writeError(w, http.StatusInternalServerError, "failed to list api keys")
			return
		}
		writeJSON(w, http.StatusOK, keys)
	}
}

func revokeAPIKeyHandler(store *auth.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID, ok := auth.TenantIDFromContext(r.Context())
		if !ok {
			writeError(w, http.StatusUnauthorized, "tenant context missing")
			return
		}
		keyID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil || keyID <= 0 {
			writeError(w, http.StatusBadRequest, "invalid key id")
			return
		}
		if err := store.RevokeAPIKey(r.Context(), tenantID, keyID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeError(w, http.StatusNotFound, "api key not found")
				return
			}
			writeError(w, http.StatusInternalServerError, "failed to revoke api key")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "revoked"})
	}
}
