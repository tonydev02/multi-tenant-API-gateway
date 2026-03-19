package gatewayhttp

import (
	"net/http"

	"github.com/namta/multi-tenant-api-gateway/backend/internal/auth"
	"github.com/namta/multi-tenant-api-gateway/backend/internal/metrics"
)

func trafficSummaryHandler(service *metrics.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID, ok := auth.TenantIDFromContext(r.Context())
		if !ok || tenantID <= 0 {
			writeError(w, http.StatusUnauthorized, "tenant context missing")
			return
		}

		writeJSON(w, http.StatusOK, service.Summary(tenantID))
	}
}
