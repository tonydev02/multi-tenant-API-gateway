# PHASE-UAT: 04 Proxy and Logging

## UAT checklist

### Proxy safety
- [ ] Proxy route requires authenticated tenant context.
- [ ] Upstream target resolution does not trust client tenant/upstream input.
- [ ] Tenant A cannot reach Tenant B upstream route mapping.

### Proxy behavior
- [ ] Valid proxy request returns upstream response.
- [ ] Upstream unavailable/timeout returns expected error status.
- [ ] Request ID is forwarded to upstream and returned to client.

### Structured logging
- [ ] Success logs include `tenant_id`, `route`, `status`, `latency_ms`, `request_id`.
- [ ] Failure logs include error context and request ID.
- [ ] Log format is consistent JSON per request.

### Verification
- [ ] `go test ./...` passes.
- [ ] `go vet ./...` passes.
- [ ] Backend build passes.
- [ ] Frontend build passes.

## Exit criteria
Phase 04 is complete when tenant-safe proxying and required request logging fields are verified in tests and runtime checks.
