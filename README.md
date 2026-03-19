# Multi-Tenant API Gateway SaaS

Production-style API gateway MVP built to demonstrate why Go is a strong fit for concurrent, network-heavy backend systems.

## Why This Project
- One shared gateway serves multiple tenants safely.
- Tenant identity is resolved server-side from trusted credentials, not client input.
- Core gateway concerns are implemented end-to-end.
- Authentication and tenancy.
- Tenant-aware rate limiting.
- Tenant-safe reverse proxying.
- Structured request logging.
- Admin dashboard workflows for tenant operations.

## What Is Implemented (Through Phase 05)
- Backend: Go (`net/http`) REST service
- Frontend: React + TypeScript admin dashboard
- Data: PostgreSQL (system of record) + Redis (rate limiting)
- Local infra: Docker Compose

### Admin capabilities
- Register tenant + admin user
- Login with JWT
- View/update current tenant
- Create/list/revoke API keys
- View tenant-scoped traffic summary

### Consumer capabilities
- Authenticate with `X-API-Key`
- Resolve tenant identity (`/api/consumer/whoami`)
- Proxy requests to tenant-scoped upstreams (`/api/consumer/proxy/...`)

## Multi-Tenancy Guarantees
- Tenant context is inferred from JWT claims (admin) or API key lookup (consumer).
- Client-supplied tenant identifiers are not trusted for routing or data access.
- Rate limiting keys are tenant-aware.
- Proxy upstream resolution is tenant + service scoped.
- Traffic summary is tenant-scoped from auth context.

## Repository Layout
- `backend/` Go API service
- `frontend/` React admin dashboard
- `docs/` architecture and API docs
- `.planning/` phase plans, UAT, and status tracking
- `docker-compose.yml` PostgreSQL + Redis for local development

## Prerequisites
- Go `1.25+`
- Node.js `20+`
- npm `10+`
- Docker + Docker Compose

## Quickstart
```bash
cp .env.example .env
set -a; source .env; set +a
make compose-up
make backend-run
```

In another terminal:
```bash
make frontend-install
cd frontend && npm run dev
```

Open `http://localhost:5173`.

Notes:
- Default local PostgreSQL port: `55432`
- Default local Redis port: `56379`
- If CORS blocks frontend calls, confirm `.env` `FRONTEND_ORIGIN` matches frontend URL.

## Verification Commands
```bash
make backend-test
make backend-vet
make backend-build
make frontend-build
make compose-config
```

## API Surface (Current)
- `GET /health`
- `POST /api/admin/tenants/register`
- `POST /api/admin/login`
- `GET /api/admin/me`
- `GET /api/admin/tenants/current`
- `PATCH /api/admin/tenants/current`
- `DELETE /api/admin/tenants/current`
- `POST /api/admin/api-keys`
- `GET /api/admin/api-keys`
- `POST /api/admin/api-keys/{id}/revoke`
- `GET /api/admin/traffic/summary`
- `GET /api/consumer/whoami`
- `ANY /api/consumer/proxy/{service}/{path...}`

## Interview Demo Script (Suggested)
1. Start infra, backend, and frontend.
2. Register/login via admin dashboard.
3. Update tenant name and show persistence.
4. Create API key and show one-time plaintext secret.
5. Use key on `/api/consumer/whoami`.
6. Revoke key and show access is denied.
7. Trigger rate limiting and show `429`.
8. Show traffic summary updates (`rate_limited_requests`, status counters).
9. Attempt tenant spoofing header and show tenant context remains auth-bound.

## Architecture and API Docs
- `docs/architecture.md`
- `docs/api-overview.md`

## Dependencies and Rationale
### Backend
- `github.com/lib/pq` PostgreSQL driver for `database/sql`.
- `golang.org/x/crypto/bcrypt` secure password hashing/verification.
- `github.com/redis/go-redis/v9` Redis client for rate-limiting counters.

### Frontend
- `react`, `react-dom`, `typescript`, `vite`, `@vitejs/plugin-react`.

### Infrastructure
- `postgres:16-alpine`
- `redis:7-alpine`

## Config Highlights
- `JWT_SECRET`, `JWT_ISSUER`, `JWT_EXPIRY_MINUTES`
- `DATABASE_URL`
- `REDIS_ADDR`, `REDIS_PASSWORD`, `REDIS_DB`
- `RATE_LIMIT_REQUESTS`, `RATE_LIMIT_WINDOW_SECONDS`
- `PROXY_TIMEOUT_SECONDS`, `PROXY_UPSTREAMS`

## Logging
- Structured JSON request logs via `log/slog`
- Required fields: `tenant_id`, `route`, `status`, `latency_ms`, `request_id`
- `X-Request-ID` is returned in responses and forwarded upstream
