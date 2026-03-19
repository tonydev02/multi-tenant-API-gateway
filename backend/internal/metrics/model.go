package metrics

// TrafficSummary is the dashboard response payload for tenant traffic visibility.
type TrafficSummary struct {
	TenantID            int64 `json:"tenant_id"`
	TotalRequests       int64 `json:"total_requests"`
	RateLimitedRequests int64 `json:"rate_limited_requests"`
	Status2xx           int64 `json:"status_2xx"`
	Status4xx           int64 `json:"status_4xx"`
	Status5xx           int64 `json:"status_5xx"`
	AvgLatencyMS        int64 `json:"avg_latency_ms"`
}
