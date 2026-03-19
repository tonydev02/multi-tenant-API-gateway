package metrics

import (
	"sync"
	"time"
)

type tenantTotals struct {
	totalRequests       int64
	rateLimitedRequests int64
	status2xx           int64
	status4xx           int64
	status5xx           int64
	totalLatencyMS      int64
}

// Service stores in-memory tenant traffic aggregates for dashboard visibility.
type Service struct {
	mu      sync.RWMutex
	tenants map[int64]tenantTotals
}

func NewService() *Service {
	return &Service{
		tenants: make(map[int64]tenantTotals),
	}
}

func (s *Service) Record(tenantID int64, statusCode int, latency time.Duration) {
	if s == nil || tenantID <= 0 {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	totals := s.tenants[tenantID]
	totals.totalRequests++
	totals.totalLatencyMS += latency.Milliseconds()

	switch {
	case statusCode == 429:
		totals.rateLimitedRequests++
		totals.status4xx++
	case statusCode >= 200 && statusCode < 300:
		totals.status2xx++
	case statusCode >= 400 && statusCode < 500:
		totals.status4xx++
	case statusCode >= 500 && statusCode < 600:
		totals.status5xx++
	}

	s.tenants[tenantID] = totals
}

func (s *Service) Summary(tenantID int64) TrafficSummary {
	if s == nil || tenantID <= 0 {
		return TrafficSummary{}
	}

	s.mu.RLock()
	totals := s.tenants[tenantID]
	s.mu.RUnlock()

	avgLatency := int64(0)
	if totals.totalRequests > 0 {
		avgLatency = totals.totalLatencyMS / totals.totalRequests
	}

	return TrafficSummary{
		TenantID:            tenantID,
		TotalRequests:       totals.totalRequests,
		RateLimitedRequests: totals.rateLimitedRequests,
		Status2xx:           totals.status2xx,
		Status4xx:           totals.status4xx,
		Status5xx:           totals.status5xx,
		AvgLatencyMS:        avgLatency,
	}
}
