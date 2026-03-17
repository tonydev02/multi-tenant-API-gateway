package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/namta/multi-tenant-api-gateway/backend/internal/auth"
)

type fakeResolver struct {
	url url.URL
	err error
}

func (f *fakeResolver) Resolve(int64, string) (url.URL, error) {
	if f.err != nil {
		return url.URL{}, f.err
	}
	return f.url, nil
}

func TestProxyHandlerForwardsRequest(t *testing.T) {
	var gotPath, gotQuery, gotMethod, gotRequestID, gotBody string
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		gotMethod = r.Method
		gotRequestID = r.Header.Get(RequestIDHeader)
		body, _ := io.ReadAll(r.Body)
		gotBody = string(body)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer upstream.Close()

	target, err := url.Parse(upstream.URL + "/base")
	if err != nil {
		t.Fatalf("parse upstream: %v", err)
	}

	h := NewHandler(&fakeResolver{url: *target}, 2*time.Second)

	req := httptest.NewRequest(http.MethodPost, "/api/consumer/proxy/billing/invoices?id=42", strings.NewReader("payload"))
	req.Header.Set(RequestIDHeader, "rid-123")
	req = req.WithContext(auth.WithClaims(req.Context(), auth.Claims{TenantID: 1}))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
	}
	if gotPath != "/base/invoices" {
		t.Fatalf("upstream path = %q", gotPath)
	}
	if gotQuery != "id=42" {
		t.Fatalf("upstream query = %q", gotQuery)
	}
	if gotMethod != http.MethodPost {
		t.Fatalf("upstream method = %q", gotMethod)
	}
	if gotRequestID != "rid-123" {
		t.Fatalf("upstream request id = %q", gotRequestID)
	}
	if gotBody != "payload" {
		t.Fatalf("upstream body = %q", gotBody)
	}
}

func TestProxyHandlerNotFound(t *testing.T) {
	h := NewHandler(&fakeResolver{err: ErrUpstreamNotFound}, time.Second)
	req := httptest.NewRequest(http.MethodGet, "/api/consumer/proxy/billing", nil)
	req = req.WithContext(auth.WithClaims(req.Context(), auth.Claims{TenantID: 1}))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestProxyHandlerBadGateway(t *testing.T) {
	target, err := url.Parse("http://127.0.0.1:1")
	if err != nil {
		t.Fatalf("parse target: %v", err)
	}

	h := NewHandler(&fakeResolver{url: *target}, 500*time.Millisecond)
	req := httptest.NewRequest(http.MethodGet, "/api/consumer/proxy/billing", nil)
	req = req.WithContext(auth.WithClaims(req.Context(), auth.Claims{TenantID: 1}))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadGateway)
	}
}
