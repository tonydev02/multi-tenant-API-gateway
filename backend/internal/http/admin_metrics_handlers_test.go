package gatewayhttp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/namta/multi-tenant-api-gateway/backend/internal/auth"
	"github.com/namta/multi-tenant-api-gateway/backend/internal/metrics"
)

func TestTrafficSummaryHandlerRequiresTenantContext(t *testing.T) {
	h := trafficSummaryHandler(metrics.NewService())

	req := httptest.NewRequest(http.MethodGet, "/api/admin/traffic/summary", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestTrafficSummaryHandlerReturnsTenantScopedSummary(t *testing.T) {
	svc := metrics.NewService()
	svc.Record(11, 200, 10*time.Millisecond)
	svc.Record(11, 429, 20*time.Millisecond)
	svc.Record(22, 500, 30*time.Millisecond)

	h := trafficSummaryHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/admin/traffic/summary", nil)
	req = req.WithContext(auth.WithClaims(req.Context(), auth.Claims{TenantID: 11}))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var payload metrics.TrafficSummary
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("parse json: %v", err)
	}
	if payload.TenantID != 11 {
		t.Fatalf("tenant_id = %d, want 11", payload.TenantID)
	}
	if payload.TotalRequests != 2 {
		t.Fatalf("total_requests = %d, want 2", payload.TotalRequests)
	}
	if payload.RateLimitedRequests != 1 {
		t.Fatalf("rate_limited_requests = %d, want 1", payload.RateLimitedRequests)
	}
	if payload.Status5xx != 0 {
		t.Fatalf("status_5xx = %d, want 0", payload.Status5xx)
	}
}
