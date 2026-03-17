package proxy

import (
	"errors"
	"testing"
)

func TestServiceResolve(t *testing.T) {
	store, err := NewMemoryStoreFromConfig("1:billing=http://localhost:18081,2:billing=http://localhost:28081")
	if err != nil {
		t.Fatalf("parse config: %v", err)
	}

	svc := NewService(store)
	tests := []struct {
		name      string
		tenantID  int64
		service   string
		wantHost  string
		wantError bool
	}{
		{name: "resolves tenant service", tenantID: 1, service: "billing", wantHost: "localhost:18081"},
		{name: "resolves case-insensitive service", tenantID: 2, service: "BILLING", wantHost: "localhost:28081"},
		{name: "missing mapping", tenantID: 3, service: "billing", wantError: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := svc.Resolve(tc.tenantID, tc.service)
			if tc.wantError {
				if err == nil {
					t.Fatal("expected error")
				}
				if !errors.Is(err, ErrUpstreamNotFound) {
					t.Fatalf("error = %v, want ErrUpstreamNotFound", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("resolve error: %v", err)
			}
			if got.Host != tc.wantHost {
				t.Fatalf("host = %q, want %q", got.Host, tc.wantHost)
			}
		})
	}
}

func TestNewMemoryStoreFromConfigInvalidEntry(t *testing.T) {
	_, err := NewMemoryStoreFromConfig("oops")
	if err == nil {
		t.Fatal("expected parse error")
	}
}
