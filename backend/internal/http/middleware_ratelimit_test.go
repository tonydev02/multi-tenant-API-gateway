package gatewayhttp

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/namta/multi-tenant-api-gateway/backend/internal/auth"
	"github.com/namta/multi-tenant-api-gateway/backend/internal/ratelimit"
)

type fakeLimiter struct {
	decision ratelimit.Decision
	err      error
}

func (f *fakeLimiter) Allow(context.Context, int64, string, ratelimit.Policy) (ratelimit.Decision, error) {
	if f.err != nil {
		return ratelimit.Decision{}, f.err
	}
	return f.decision, nil
}

func TestRateLimitMiddlewareAllows(t *testing.T) {
	mw := requireTenantRateLimit(&fakeLimiter{decision: ratelimit.Decision{Allowed: true}}, ratelimit.Policy{Requests: 10, Window: time.Minute})
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/admin/me", nil)
	req = req.WithContext(auth.WithClaims(req.Context(), auth.Claims{TenantID: 1}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestRateLimitMiddlewareBlocks(t *testing.T) {
	mw := requireTenantRateLimit(&fakeLimiter{decision: ratelimit.Decision{Allowed: false, Limit: 1, Remaining: 0, ResetAt: time.Now().Add(time.Minute)}}, ratelimit.Policy{Requests: 1, Window: time.Minute})
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/admin/me", nil)
	req = req.WithContext(auth.WithClaims(req.Context(), auth.Claims{TenantID: 1}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusTooManyRequests)
	}
}

func TestRateLimitMiddlewareErrors(t *testing.T) {
	mw := requireTenantRateLimit(&fakeLimiter{err: errors.New("redis down")}, ratelimit.Policy{Requests: 10, Window: time.Minute})
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/admin/me", nil)
	req = req.WithContext(auth.WithClaims(req.Context(), auth.Claims{TenantID: 1}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
}

func TestNormalizedRoute(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		want   string
	}{
		{
			name:   "numeric path id",
			method: http.MethodGet,
			path:   "/api/admin/api-keys/42/revoke",
			want:   "GET:/api/admin/api-keys/:id/revoke",
		},
		{
			name:   "uuid path id",
			method: http.MethodPost,
			path:   "/v1/orders/123e4567-e89b-12d3-a456-426614174000/items",
			want:   "POST:/v1/orders/:id/items",
		},
		{
			name:   "ulid path id",
			method: http.MethodDelete,
			path:   "/v1/sessions/01ARZ3NDEKTSV4RRFFQ69G5FAV",
			want:   "DELETE:/v1/sessions/:id",
		},
		{
			name:   "static path",
			method: http.MethodPatch,
			path:   "/api/admin/tenants/current",
			want:   "PATCH:/api/admin/tenants/current",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			got := normalizedRoute(req)
			if got != tc.want {
				t.Fatalf("normalizedRoute() = %q, want %q", got, tc.want)
			}
		})
	}
}
