package proxy

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// MemoryStore stores tenant/service upstream mappings in memory.
type MemoryStore struct {
	upstreams map[UpstreamKey]url.URL
}

// NewMemoryStoreFromConfig parses PROXY_UPSTREAMS format:
// <tenant_id>:<service>=<base_url>,<tenant_id>:<service>=<base_url>
func NewMemoryStoreFromConfig(raw string) (*MemoryStore, error) {
	store := &MemoryStore{upstreams: make(map[UpstreamKey]url.URL)}
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return store, nil
	}

	parts := strings.Split(trimmed, ",")
	for _, part := range parts {
		entry := strings.TrimSpace(part)
		if entry == "" {
			continue
		}

		mapping, err := parseEntry(entry)
		if err != nil {
			return nil, err
		}
		store.upstreams[mapping.Key] = mapping.BaseURL
	}

	return store, nil
}

func (s *MemoryStore) Get(tenantID int64, service string) (url.URL, bool) {
	u, ok := s.upstreams[UpstreamKey{TenantID: tenantID, Service: strings.ToLower(service)}]
	return u, ok
}

func parseEntry(entry string) (UpstreamMapping, error) {
	left, right, ok := strings.Cut(entry, "=")
	if !ok {
		return UpstreamMapping{}, fmt.Errorf("invalid upstream entry %q: expected '='", entry)
	}

	tenantRaw, service, ok := strings.Cut(strings.TrimSpace(left), ":")
	if !ok {
		return UpstreamMapping{}, fmt.Errorf("invalid upstream entry %q: expected tenant:service", entry)
	}

	tenantID, err := strconv.ParseInt(strings.TrimSpace(tenantRaw), 10, 64)
	if err != nil || tenantID <= 0 {
		return UpstreamMapping{}, fmt.Errorf("invalid upstream entry %q: invalid tenant id", entry)
	}

	service = strings.ToLower(strings.TrimSpace(service))
	if service == "" {
		return UpstreamMapping{}, fmt.Errorf("invalid upstream entry %q: service is required", entry)
	}

	baseURL, err := url.Parse(strings.TrimSpace(right))
	if err != nil {
		return UpstreamMapping{}, fmt.Errorf("invalid upstream entry %q: parse url: %w", entry, err)
	}
	if baseURL.Scheme == "" || baseURL.Host == "" {
		return UpstreamMapping{}, fmt.Errorf("invalid upstream entry %q: url must include scheme and host", entry)
	}

	return UpstreamMapping{
		Key: UpstreamKey{
			TenantID: tenantID,
			Service:  service,
		},
		BaseURL: *baseURL,
	}, nil
}
