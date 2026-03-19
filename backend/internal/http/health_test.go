package gatewayhttp

import (
	"context"
	"errors"
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

func TestReadyHandler(t *testing.T) {
	tests := []struct {
		name           string
		check          func() func(ctx context.Context) error
		expectStatus   int
		expectContains string
	}{
		{
			name: "ready when dependencies are healthy",
			check: func() func(ctx context.Context) error {
				return func(ctx context.Context) error { return nil }
			},
			expectStatus:   http.StatusOK,
			expectContains: `"status":"ok"`,
		},
		{
			name: "degraded when dependency check fails",
			check: func() func(ctx context.Context) error {
				return func(ctx context.Context) error { return errors.New("boom") }
			},
			expectStatus:   http.StatusServiceUnavailable,
			expectContains: `"status":"degraded"`,
		},
		{
			name:           "degraded when checker is missing",
			check:          func() func(ctx context.Context) error { return nil },
			expectStatus:   http.StatusServiceUnavailable,
			expectContains: `"status":"degraded"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
			recorder := httptest.NewRecorder()

			readyHandler(tc.check())(recorder, req)

			if recorder.Code != tc.expectStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, tc.expectStatus)
			}
			if ct := recorder.Header().Get("Content-Type"); ct != "application/json" {
				t.Fatalf("content-type = %q, want application/json", ct)
			}
			if !strings.Contains(recorder.Body.String(), tc.expectContains) {
				t.Fatalf("body = %q, want substring %q", recorder.Body.String(), tc.expectContains)
			}
		})
	}
}
