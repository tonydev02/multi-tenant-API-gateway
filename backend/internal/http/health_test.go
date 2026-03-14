package gatewayhttp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		expectStatus int
	}{
		{name: "get health", method: http.MethodGet, expectStatus: http.StatusOK},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/health", nil)
			recorder := httptest.NewRecorder()

			healthHandler(recorder, req)

			if recorder.Code != tc.expectStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, tc.expectStatus)
			}
			if ct := recorder.Header().Get("Content-Type"); ct != "application/json" {
				t.Fatalf("content-type = %q, want application/json", ct)
			}
			if !strings.Contains(recorder.Body.String(), `"status":"ok"`) {
				t.Fatalf("body = %q, want status ok payload", recorder.Body.String())
			}
		})
	}
}
