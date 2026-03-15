package ratelimit

import (
	"context"
	"fmt"
	"time"
)

// CounterStore tracks counters scoped by tenant, route, and window.
type CounterStore interface {
	Increment(ctx context.Context, tenantID int64, route string, windowStart time.Time, window time.Duration) (int64, error)
}

// Service evaluates policies using a counter store.
type Service struct {
	store CounterStore
	nowFn func() time.Time
}

func NewService(store CounterStore) *Service {
	return &Service{store: store, nowFn: time.Now}
}

func (s *Service) Allow(ctx context.Context, tenantID int64, route string, policy Policy) (Decision, error) {
	if policy.Requests <= 0 || policy.Window <= 0 {
		return Decision{}, ErrInvalidPolicy
	}
	if tenantID <= 0 || route == "" {
		return Decision{}, ErrInvalidPolicy
	}

	now := s.nowFn().UTC()
	windowStart := now.Truncate(policy.Window)
	count, err := s.store.Increment(ctx, tenantID, route, windowStart, policy.Window)
	if err != nil {
		return Decision{}, fmt.Errorf("increment rate limit counter: %w", err)
	}

	remaining := policy.Requests - count
	if remaining < 0 {
		remaining = 0
	}

	return Decision{
		Allowed:   count <= policy.Requests,
		Limit:     policy.Requests,
		Remaining: remaining,
		ResetAt:   windowStart.Add(policy.Window),
	}, nil
}
