package ratelimit

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
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

func NewRedisClient(ctx context.Context, addr, username, password string, db int, useTLS bool) (*redis.Client, error) {
	var (
		client *redis.Client
		err    error
	)

	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(addr)), "redis://") || strings.HasPrefix(strings.ToLower(strings.TrimSpace(addr)), "rediss://") {
		client, err = newRedisClientFromURL(addr, username, password, db)
		if err != nil {
			return nil, err
		}
	} else {
		client = newRedisClientFromAddr(addr, username, password, db, useTLS)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return client, nil
}

func newRedisClientFromURL(rawURL, username, password string, db int) (*redis.Client, error) {
	opts, err := redis.ParseURL(rawURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}
	if strings.TrimSpace(username) != "" {
		opts.Username = username
	}
	if strings.TrimSpace(password) != "" {
		opts.Password = password
	}
	opts.DB = db
	return redis.NewClient(opts), nil
}

func newRedisClientFromAddr(addr, username, password string, db int, useTLS bool) *redis.Client {
	opts := &redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       db,
	}
	if useTLS {
		serverName := addr
		if host, _, err := net.SplitHostPort(addr); err == nil && host != "" {
			serverName = host
		}
		opts.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: serverName,
		}
	}
	return redis.NewClient(opts)
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
