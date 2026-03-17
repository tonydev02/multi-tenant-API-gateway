package proxy

import "net/url"

// UpstreamKey identifies a tenant-specific service mapping.
type UpstreamKey struct {
	TenantID int64
	Service  string
}

// UpstreamMapping stores the upstream base URL for a tenant/service.
type UpstreamMapping struct {
	Key     UpstreamKey
	BaseURL url.URL
}
