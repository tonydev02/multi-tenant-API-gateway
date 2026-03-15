package ratelimit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

type fakeRedisClient struct {
	lastIncrKey   string
	lastExpireKey string
	lastExpireTTL time.Duration
	value         int64
	err           error
}

func (f *fakeRedisClient) Incr(ctx context.Context, key string) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx)
	if f.err != nil {
		cmd.SetErr(f.err)
		return cmd
	}
	f.lastIncrKey = key
	f.value++
	cmd.SetVal(f.value)
	return cmd
}

func (f *fakeRedisClient) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	cmd := redis.NewBoolCmd(ctx)
	if f.err != nil {
		cmd.SetErr(f.err)
		return cmd
	}
	f.lastExpireKey = key
	f.lastExpireTTL = expiration
	cmd.SetVal(true)
	return cmd
}

func TestRedisStoreIncrementSetsExpiryOnFirstIncrement(t *testing.T) {
	fake := &fakeRedisClient{}
	store := NewRedisStore(fake)

	count, err := store.Increment(context.Background(), 12, "GET:/api/admin/me", time.Unix(120, 0), time.Minute)
	if err != nil {
		t.Fatalf("increment: %v", err)
	}
	if count != 1 {
		t.Fatalf("count = %d, want 1", count)
	}
	if fake.lastIncrKey == "" || fake.lastExpireKey == "" {
		t.Fatalf("expected key + expiry to be recorded")
	}
	if fake.lastExpireTTL <= time.Minute {
		t.Fatalf("expiry ttl = %s, want > 1m", fake.lastExpireTTL)
	}
}

func TestRedisStoreIncrementReturnsError(t *testing.T) {
	fake := &fakeRedisClient{err: errors.New("boom")}
	store := NewRedisStore(fake)

	if _, err := store.Increment(context.Background(), 1, "GET:/api/admin/me", time.Unix(0, 0), time.Minute); err == nil {
		t.Fatal("expected error")
	}
}
