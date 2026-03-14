package gatewayhttp

import (
	"net/http"

	"github.com/namta/multi-tenant-api-gateway/backend/internal/auth"
)

func requireTenantContext() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if tenantID, ok := auth.TenantIDFromContext(r.Context()); !ok || tenantID <= 0 {
				writeError(w, http.StatusUnauthorized, "tenant context missing")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
