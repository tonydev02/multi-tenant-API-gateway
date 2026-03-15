package ratelimit

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type fakeCounterStore struct {
	counts map[string]int64
}

func (f *fakeCounterStore) Increment(_ context.Context, tenantID int64, route string, windowStart time.Time, _ time.Duration) (int64, error) {
	if f.counts == nil {
		f.counts = map[string]int64{}
	}
	key := fmt.Sprintf("%d|%s|%d", tenantID, route, windowStart.Unix())
	f.counts[key]++
	return f.counts[key], nil
}

func TestAllowDecision(t *testing.T) {
	store := &fakeCounterStore{}
	svc := NewService(store)
	svc.nowFn = func() time.Time { return time.Unix(100, 0) }

	policy := Policy{Requests: 2, Window: time.Minute}

	first, err := svc.Allow(context.Background(), 1, "GET:/api/admin/me", policy)
	if err != nil {
		t.Fatalf("first allow: %v", err)
	}
	if !first.Allowed || first.Remaining != 1 {
		t.Fatalf("first decision = %+v", first)
	}

	second, err := svc.Allow(context.Background(), 1, "GET:/api/admin/me", policy)
	if err != nil {
		t.Fatalf("second allow: %v", err)
	}
	if !second.Allowed || second.Remaining != 0 {
		t.Fatalf("second decision = %+v", second)
	}

	third, err := svc.Allow(context.Background(), 1, "GET:/api/admin/me", policy)
	if err != nil {
		t.Fatalf("third allow: %v", err)
	}
	if third.Allowed || third.Remaining != 0 {
		t.Fatalf("third decision = %+v", third)
	}
}

func TestAllowWindowReset(t *testing.T) {
	store := &fakeCounterStore{}
	svc := NewService(store)

	calls := 0
	svc.nowFn = func() time.Time {
		calls++
		if calls == 1 {
			return time.Unix(100, 0)
		}
		return time.Unix(161, 0)
	}

	policy := Policy{Requests: 1, Window: time.Minute}

	first, err := svc.Allow(context.Background(), 1, "GET:/api/admin/me", policy)
	if err != nil {
		t.Fatalf("first allow: %v", err)
	}
	if !first.Allowed {
		t.Fatalf("first decision = %+v", first)
	}

	second, err := svc.Allow(context.Background(), 1, "GET:/api/admin/me", policy)
	if err != nil {
		t.Fatalf("second allow: %v", err)
	}
	if !second.Allowed {
		t.Fatalf("second decision should reset window, got %+v", second)
	}
}
