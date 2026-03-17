package gatewayhttp

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/namta/multi-tenant-api-gateway/backend/internal/auth"
	"github.com/namta/multi-tenant-api-gateway/backend/internal/proxy"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

const gatewayTenantIDHeader = "X-Gateway-Tenant-ID"

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func withRequestLogging(logger *slog.Logger) func(http.Handler) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			started := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(rec, r)

			requestID, _ := requestIDFromContext(r.Context())
			tenantID := int64(0)
			if rawTenantID := r.Header.Get(gatewayTenantIDHeader); rawTenantID != "" {
				if parsed, err := strconv.ParseInt(rawTenantID, 10, 64); err == nil {
					tenantID = parsed
				}
			} else if parsed, ok := auth.TenantIDFromContext(r.Context()); ok {
				tenantID = parsed
			}

			attrs := []any{
				"request_id", requestID,
				"tenant_id", tenantID,
				"route", normalizedRoute(r),
				"status", rec.status,
				"latency_ms", time.Since(started).Milliseconds(),
				"method", r.Method,
				"path", r.URL.Path,
			}
			if upstreamHost := r.Header.Get(proxy.ProxyUpstreamHostHeader); upstreamHost != "" {
				attrs = append(attrs, "upstream_host", upstreamHost)
			}

			logger.Info("request complete", attrs...)
		})
	}
}
