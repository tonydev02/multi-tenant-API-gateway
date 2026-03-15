# PHASE-PLAN: 03 Rate Limiting

## Scope
Implement Redis-backed tenant-aware rate limiting middleware and integrate it into backend request flows.

## Files to create/edit

### Backend: rate limit core
- `backend/internal/ratelimit/model.go` (policy structs, decision result)
- `backend/internal/ratelimit/redis_store.go` (counter increment/expiry)
- `backend/internal/ratelimit/service.go` (window math + decision)
- `backend/internal/ratelimit/errors.go`

### Backend: HTTP integration
- `backend/internal/http/middleware_ratelimit.go`
- `backend/internal/http/router.go` (attach middleware to route groups)
- `backend/internal/http/response.go` (consistent 429 response shape)

### Backend: configuration
- `backend/internal/config/config.go` (rate-limit env vars)
- `.env.example` (rate-limit defaults + redis connection vars if needed)
- `backend/cmd/server/main.go` (wire Redis limiter dependencies)

### Backend: data/infra
- `docker-compose.yml` (confirm redis service config for limiter use)

### Tests
- `backend/internal/ratelimit/service_test.go` (table-driven decision tests)
- `backend/internal/ratelimit/redis_store_test.go` (integration-like behavior via test doubles)
- `backend/internal/http/middleware_ratelimit_test.go`

### Documentation
- `README.md` (rate-limit env vars and behavior)
- `docs/architecture.md` (rate limiting section)
- `.planning/STATE.md`
- `.planning/phases/03-rate-limiting/PHASE-UAT.md`
- `.planning/phases/03-rate-limiting/PHASE-SUMMARY.md`

## Acceptance criteria
- Limiting is tenant-aware and based on context tenant ID.
- Exceeding configured threshold returns HTTP 429 with clear response.
- Requests under threshold continue normally.
- Redis counter keys include tenant and route scope.
- Limiter middleware is applied to target protected routes.
- Configurable defaults via environment variables are documented.
- Backend tests cover allow/deny decisions and window reset behavior.

## Verification commands
- `cd backend && go test ./...`
- `cd backend && go vet ./...`
- `cd backend && go build ./...`
- `cd frontend && npm run build`
- `docker compose config`
- Runtime smoke checks:
  - Repeat protected endpoint requests until 429 is returned.
  - Wait for next window and verify requests are accepted again.

## Risks
- Misconfigured limits can block legitimate traffic.
- Route-level policy mapping can drift as routes evolve.
- Redis latency may increase request latency.

## Non-goals
- Admin UI for dynamic limit edits.
- Per-user or per-IP secondary limits.
- Billing-tier policy engine.
