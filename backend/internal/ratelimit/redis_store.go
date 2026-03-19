package ratelimit

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCounterClient interface {
	Incr(ctx context.Context, key string) *redis.IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
}

// RedisStore is the counter store backed by Redis.
type RedisStore struct {
	client redisCounterClient
}

func NewRedisStore(client redisCounterClient) *RedisStore {
	return &RedisStore{client: client}
}

func NewRedisClient(ctx context.Context, addr, password string, db int, useTLS bool) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	}
	if useTLS {
		opts.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	client := redis.NewClient(opts)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return client, nil
}

func (s *RedisStore) Increment(ctx context.Context, tenantID int64, route string, windowStart time.Time, window time.Duration) (int64, error) {
	key := rateLimitKey(tenantID, route, windowStart)
	count, err := s.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if count == 1 {
		if err := s.client.Expire(ctx, key, window+5*time.Second).Err(); err != nil {
			return 0, err
		}
	}

	return count, nil
}

func rateLimitKey(tenantID int64, route string, windowStart time.Time) string {
	safeRoute := strings.NewReplacer(" ", "_", "/", ":").Replace(route)
	return fmt.Sprintf("rl:%d:%s:%d", tenantID, safeRoute, windowStart.Unix())
}
