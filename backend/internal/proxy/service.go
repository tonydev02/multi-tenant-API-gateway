package proxy

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var ErrUpstreamNotFound = errors.New("upstream not found")

// Resolver resolves an upstream for a tenant and service.
type Resolver interface {
	Resolve(tenantID int64, service string) (url.URL, error)
}

// Store describes a source of tenant/service upstream mappings.
type Store interface {
	Get(tenantID int64, service string) (url.URL, bool)
}

// Service resolves tenant-safe upstream mappings.
type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) Resolve(tenantID int64, service string) (url.URL, error) {
	if tenantID <= 0 {
		return url.URL{}, fmt.Errorf("resolve upstream: tenant id must be positive")
	}
	normalizedService := strings.ToLower(strings.TrimSpace(service))
	if normalizedService == "" {
		return url.URL{}, fmt.Errorf("resolve upstream: service is required")
	}
	if s.store == nil {
		return url.URL{}, ErrUpstreamNotFound
	}

	u, ok := s.store.Get(tenantID, normalizedService)
	if !ok {
		return url.URL{}, ErrUpstreamNotFound
	}
	return u, nil
}
