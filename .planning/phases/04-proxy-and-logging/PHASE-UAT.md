# PHASE-UAT: 04 Proxy and Logging

## UAT checklist

### Proxy safety
- [x] Proxy route requires authenticated tenant context.
- [x] Upstream target resolution does not trust client tenant/upstream input.
- [x] Tenant A cannot reach Tenant B upstream route mapping.

### Proxy behavior
- [x] Valid proxy request returns upstream response.
- [x] Upstream unavailable/timeout returns expected error status.
- [x] Request ID is forwarded to upstream and returned to client.

### Structured logging
- [x] Success logs include `tenant_id`, `route`, `status`, `latency_ms`, `request_id`.
- [x] Failure logs include error context and request ID.
- [x] Log format is consistent JSON per request.

### Verification
- [x] `go test ./...` passes.
- [x] `go vet ./...` passes.
- [x] Backend build passes.
- [x] Frontend build passes.

## Exit criteria
Phase 04 is complete when tenant-safe proxying and required request logging fields are verified in tests and runtime checks.
