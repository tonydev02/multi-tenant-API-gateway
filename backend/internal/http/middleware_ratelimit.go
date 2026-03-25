package gatewayhttp

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/namta/multi-tenant-api-gateway/backend/internal/auth"
	"github.com/namta/multi-tenant-api-gateway/backend/internal/ratelimit"
)

type tenantRateLimiter interface {
	Allow(ctx context.Context, tenantID int64, route string, policy ratelimit.Policy) (ratelimit.Decision, error)
}

var (
	numericSegmentRE = regexp.MustCompile(`^\d+$`)
	uuidSegmentRE    = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
	ulidSegmentRE    = regexp.MustCompile(`^[0-9A-HJKMNP-TV-Z]{26}$`)
)

func requireTenantRateLimit(limiter tenantRateLimiter, policy ratelimit.Policy) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenantID, ok := auth.TenantIDFromContext(r.Context())
			if !ok || tenantID <= 0 {
				writeError(w, http.StatusUnauthorized, "tenant context missing")
				return
			}
			if limiter == nil {
				next.ServeHTTP(w, r)
				return
			}

			route := normalizedRoute(r)
			decision, err := limiter.Allow(r.Context(), tenantID, route, policy)
			if err != nil {
				status := http.StatusServiceUnavailable
				if errors.Is(err, ratelimit.ErrInvalidPolicy) {
					status = http.StatusInternalServerError
				}
				writeError(w, status, "rate limiter unavailable")
				return
			}
			if !decision.Allowed {
				writeJSON(w, http.StatusTooManyRequests, map[string]any{
					"error":     "rate limit exceeded",
					"limit":     decision.Limit,
					"remaining": decision.Remaining,
					"reset_at":  decision.ResetAt.UTC().Format(time.RFC3339),
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func normalizedRoute(r *http.Request) string {
	// Replace variable identifiers to avoid high-cardinality route keys.
	segments := strings.Split(r.URL.Path, "/")
	for i := range segments {
		if isVariablePathSegment(segments[i]) {
			segments[i] = ":id"
		}
	}

	return r.Method + ":" + strings.Join(segments, "/")
}

func isVariablePathSegment(segment string) bool {
	if segment == "" {
		return false
	}
	return numericSegmentRE.MatchString(segment) || uuidSegmentRE.MatchString(segment) || ulidSegmentRE.MatchString(segment)
}
