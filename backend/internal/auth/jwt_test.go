package auth

import (
	"testing"
	"time"
)

func TestJWTManagerIssueAndParse(t *testing.T) {
	manager := NewJWTManager("test-secret", "test-issuer", time.Minute)
	manager.nowFn = func() time.Time {
		return time.Unix(1000, 0)
	}

	token, err := manager.Issue(AdminUser{ID: 42, TenantID: 7, Email: "admin@example.com"})
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}

	claims, err := manager.Parse(token)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	if claims.Subject != 42 || claims.TenantID != 7 || claims.Email != "admin@example.com" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestJWTManagerRejectsExpiredToken(t *testing.T) {
	manager := NewJWTManager("test-secret", "test-issuer", time.Minute)
	manager.nowFn = func() time.Time {
		return time.Unix(1000, 0)
	}

	token, err := manager.Issue(AdminUser{ID: 1, TenantID: 1, Email: "admin@example.com"})
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}

	manager.nowFn = func() time.Time {
		return time.Unix(2000, 0)
	}
	if _, err := manager.Parse(token); err == nil {
		t.Fatal("expected parse to fail for expired token")
	}
}
