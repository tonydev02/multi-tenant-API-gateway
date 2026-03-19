package gatewayhttp

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"strings"
	"testing"
	"time"

	"github.com/namta/multi-tenant-api-gateway/backend/internal/auth"
	"github.com/namta/multi-tenant-api-gateway/backend/internal/proxy"
)

type stubTrafficRecorder struct {
	mu       sync.Mutex
	recorded []struct {
		tenantID int64
		status   int
		latency  time.Duration
	}
}

func (s *stubTrafficRecorder) Record(tenantID int64, statusCode int, latency time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.recorded = append(s.recorded, struct {
		tenantID int64
		status   int
		latency  time.Duration
	}{
		tenantID: tenantID,
		status:   statusCode,
		latency:  latency,
	})
}

func TestLoggingMiddlewareIncludesRequiredFields(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{}))
	recorder := &stubTrafficRecorder{}

	h := withRequestLogging(logger, recorder)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set(proxy.ProxyUpstreamHostHeader, "upstream.local:8080")
		w.WriteHeader(http.StatusCreated)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/consumer/proxy/billing/invoices", nil)
	req = req.WithContext(auth.WithClaims(req.Context(), auth.Claims{TenantID: 9}))
	req = req.WithContext(withRequestIDContext(req.Context(), "rid-xyz"))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	line := strings.TrimSpace(buf.String())
	if line == "" {
		t.Fatal("expected log output")
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(line), &payload); err != nil {
		t.Fatalf("parse log json: %v", err)
	}

	for _, key := range []string{"request_id", "tenant_id", "route", "status", "latency_ms"} {
		if _, ok := payload[key]; !ok {
			t.Fatalf("missing key %q in log payload: %v", key, payload)
		}
	}
	if payload["request_id"] != "rid-xyz" {
		t.Fatalf("request_id = %v", payload["request_id"])
	}
	if int(payload["status"].(float64)) != http.StatusCreated {
		t.Fatalf("status = %v", payload["status"])
	}
	if payload["upstream_host"] != "upstream.local:8080" {
		t.Fatalf("upstream_host = %v", payload["upstream_host"])
	}
	if len(recorder.recorded) != 1 {
		t.Fatalf("recorded count = %d, want 1", len(recorder.recorded))
	}
	if recorder.recorded[0].tenantID != 9 {
		t.Fatalf("recorded tenant = %d, want 9", recorder.recorded[0].tenantID)
	}
	if recorder.recorded[0].status != http.StatusCreated {
		t.Fatalf("recorded status = %d, want %d", recorder.recorded[0].status, http.StatusCreated)
	}
}
