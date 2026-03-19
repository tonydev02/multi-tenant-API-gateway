package metrics

import (
	"testing"
	"time"
)

func TestServiceSummary(t *testing.T) {
	svc := NewService()
	svc.Record(1, 200, 120*time.Millisecond)
	svc.Record(1, 201, 80*time.Millisecond)
	svc.Record(1, 429, 50*time.Millisecond)
	svc.Record(1, 404, 20*time.Millisecond)
	svc.Record(1, 500, 30*time.Millisecond)

	got := svc.Summary(1)
	if got.TenantID != 1 {
		t.Fatalf("tenant_id = %d, want 1", got.TenantID)
	}
	if got.TotalRequests != 5 {
		t.Fatalf("total_requests = %d, want 5", got.TotalRequests)
	}
	if got.RateLimitedRequests != 1 {
		t.Fatalf("rate_limited_requests = %d, want 1", got.RateLimitedRequests)
	}
	if got.Status2xx != 2 {
		t.Fatalf("status_2xx = %d, want 2", got.Status2xx)
	}
	if got.Status4xx != 2 {
		t.Fatalf("status_4xx = %d, want 2", got.Status4xx)
	}
	if got.Status5xx != 1 {
		t.Fatalf("status_5xx = %d, want 1", got.Status5xx)
	}
	if got.AvgLatencyMS != 60 {
		t.Fatalf("avg_latency_ms = %d, want 60", got.AvgLatencyMS)
	}
}

func TestServiceIgnoresInvalidTenant(t *testing.T) {
	svc := NewService()
	svc.Record(0, 200, 10*time.Millisecond)
	svc.Record(-1, 200, 10*time.Millisecond)

	got := svc.Summary(0)
	if got.TotalRequests != 0 {
		t.Fatalf("total_requests = %d, want 0", got.TotalRequests)
	}
}
