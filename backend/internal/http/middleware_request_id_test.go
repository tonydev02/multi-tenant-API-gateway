package gatewayhttp

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestIDMiddlewareGeneratesWhenMissing(t *testing.T) {
	h := withRequestID()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := requestIDFromContext(r.Context())
		if !ok || id == "" {
			t.Fatal("request id missing from context")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Header().Get(requestIDHeader) == "" {
		t.Fatal("expected response request id header")
	}
}

func TestRequestIDMiddlewarePreservesIncomingID(t *testing.T) {
	h := withRequestID()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, _ := requestIDFromContext(r.Context())
		if id != "incoming-123" {
			t.Fatalf("context id = %q", id)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set(requestIDHeader, "incoming-123")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if got := rec.Header().Get(requestIDHeader); got != "incoming-123" {
		t.Fatalf("response id = %q", got)
	}
}
