package gatewayhttp

import (
	"net/http"
	"strings"

	"github.com/namta/multi-tenant-api-gateway/backend/internal/auth"
)

func requireAdminAuth(jwtManager *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				writeError(w, http.StatusUnauthorized, "missing bearer token")
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwtManager.Parse(token)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			next.ServeHTTP(w, r.WithContext(auth.WithClaims(r.Context(), claims)))
		})
	}
}

func requireAPIKeyAuth(apiKeyAuth *auth.APIKeyAuthenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rawKey := r.Header.Get("X-API-Key")
			if rawKey == "" {
				writeError(w, http.StatusUnauthorized, "missing api key")
				return
			}
			record, err := apiKeyAuth.Authenticate(r.Context(), rawKey)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "invalid api key")
				return
			}

			claims := auth.Claims{TenantID: record.TenantID}
			next.ServeHTTP(w, r.WithContext(auth.WithClaims(r.Context(), claims)))
		})
	}
}

func chainMiddleware(middleware ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		wrapped := final
		for i := len(middleware) - 1; i >= 0; i-- {
			wrapped = middleware[i](wrapped)
		}
		return wrapped
	}
}
