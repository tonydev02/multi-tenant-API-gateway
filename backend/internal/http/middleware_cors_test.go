package gatewayhttp

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSPreflightAllowedOrigin(t *testing.T) {
	router := NewRouter(Dependencies{FrontendOrigin: "http://localhost:5173"})

	req := httptest.NewRequest(http.MethodOptions, "/api/admin/login", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("allow-origin = %q, want %q", got, "http://localhost:5173")
	}
}

func TestCORSNoHeaderForDifferentOrigin(t *testing.T) {
	router := NewRouter(Dependencies{FrontendOrigin: "http://localhost:5173"})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "http://evil.local")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("allow-origin = %q, want empty", got)
	}
}
