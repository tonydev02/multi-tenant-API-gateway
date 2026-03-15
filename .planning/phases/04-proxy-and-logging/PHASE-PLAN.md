# PHASE-PLAN: 04 Proxy and Logging

## Scope
Implement tenant-safe reverse proxy routing and structured gateway request logging.

## Suggested subphases (optional)
- **04A Proxy Foundation**: upstream config model, tenant-safe route resolution, reverse proxy handler.
- **04B Logging & Request IDs**: request ID middleware, structured access logs, proxy error logs.

## Files to create/edit

### Backend: proxy core
- `backend/internal/proxy/model.go` (upstream config + route mapping model)
- `backend/internal/proxy/store_memory.go` (MVP in-memory upstream map)
- `backend/internal/proxy/service.go` (tenant-safe upstream resolution)
- `backend/internal/proxy/handler.go` (reverse proxy execution)

### Backend: HTTP middleware/logging
- `backend/internal/http/middleware_request_id.go`
- `backend/internal/http/middleware_logging.go`
- `backend/internal/http/middleware_proxy_authz.go` (tenant-safe proxy guard)
- `backend/internal/http/router.go` (register `/proxy/...` routes + middleware order)
- `backend/internal/http/response.go` (error response consistency)

### Backend: config
- `backend/internal/config/config.go` (proxy timeout/log level env vars)
- `.env.example` (proxy/logging-related env vars)
- `backend/cmd/server/main.go` (wire proxy + logger components)

### Tests
- `backend/internal/proxy/service_test.go`
- `backend/internal/proxy/handler_test.go`
- `backend/internal/http/middleware_logging_test.go`
- `backend/internal/http/middleware_request_id_test.go`

### Documentation
- `README.md` (proxy usage + logging behavior)
- `docs/architecture.md` (proxy/logging architecture section)
- `docs/api-overview.md` (proxy endpoints)
- `.planning/STATE.md`
- `.planning/phases/04-proxy-and-logging/PHASE-UAT.md`
- `.planning/phases/04-proxy-and-logging/PHASE-SUMMARY.md`

## Acceptance criteria
- Proxy route resolves upstream using tenant-safe server-side logic only.
- Client-supplied tenant/upstream identifiers are ignored for routing.
- Proxy can forward requests/responses to configured upstream targets.
- Every gateway request log contains `tenant_id`, `route`, `status`, `latency_ms`, and `request_id`.
- Request IDs are propagated through response headers and upstream requests.
- Error paths return correct status codes and log structured error context.

## Verification commands
- `cd backend && go test ./...`
- `cd backend && go vet ./...`
- `cd backend && go build ./...`
- `cd frontend && npm run build`
- `docker compose config`
- Runtime smoke checks:
  - Authenticated tenant request through proxy returns upstream response.
  - Invalid/missing tenant context for proxy route is rejected.
  - Logs for successful and failed proxy requests contain required fields.

## Risks
- Middleware ordering mistakes can bypass auth/tenant checks.
- Proxy timeouts not tuned can produce false failures.
- Logging too early/late can miss final status or latency.

## Non-goals
- Multi-upstream load balancing algorithms.
- Distributed tracing backends.
- Full production SIEM/log shipping integration.
