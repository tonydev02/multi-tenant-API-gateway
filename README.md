# Multi-Tenant API Gateway SaaS

Production-style multi-tenant API gateway scaffold with a Go backend, React + TypeScript admin UI, and PostgreSQL/Redis local dependencies.

## Project layout
- `backend/`: Go API service with health, auth, tenancy, and API key flows
- `frontend/`: React + TypeScript admin shell with login/tenant registration
- `.planning/`: project and phase planning docs
- `docker-compose.yml`: local PostgreSQL and Redis

## Prerequisites
- Go 1.25+
- Node.js 20+
- npm 10+
- Docker + Docker Compose

## Quickstart
1. Copy env defaults:
   - `cp .env.example .env`
2. Load environment variables:
   - `set -a; source .env; set +a`
3. Start data services:
   - `make compose-up`
4. Start backend:
   - `make backend-run`
5. In another terminal, install frontend deps and run dev server:
   - `make frontend-install`
   - `cd frontend && npm run dev`

By default this project maps Docker ports to `55432` (PostgreSQL) and `56379` (Redis) to avoid conflicts with local host services.
If frontend calls are blocked by CORS, ensure `FRONTEND_ORIGIN` in `.env` matches your frontend URL (default `http://localhost:5173`).

## Verification commands
- `make backend-test`
- `make backend-vet`
- `make backend-build`
- `make frontend-build`
- `make compose-config`
- Runtime rate-limit smoke check (example):
  - call a protected endpoint repeatedly and verify HTTP `429` after threshold

## API endpoints (current)
- `GET /health` -> `200 {"status":"ok"}`
- `POST /api/admin/tenants/register` -> create tenant + admin user
- `POST /api/admin/login` -> get JWT
- `GET /api/admin/me` -> validate JWT and return claims
- `GET/PATCH/DELETE /api/admin/tenants/current` -> tenant CRUD on current tenant
- `POST /api/admin/api-keys` -> create API key (plaintext returned once)
- `GET /api/admin/api-keys` -> list tenant API keys
- `POST /api/admin/api-keys/{id}/revoke` -> revoke key
- `GET /api/consumer/whoami` -> tenant resolution via `X-API-Key`
- `ANY /api/consumer/proxy/{service}/{path...}` -> tenant-safe upstream proxying via `X-API-Key`

## Documentation
- `docs/architecture.md`: current architecture, request flows, tenancy boundaries, and data model.
- `docs/api-overview.md`: quick endpoint map for implemented MVP APIs.

## Dependencies added and why

### Backend
- `github.com/lib/pq`: PostgreSQL driver for `database/sql`.
- `golang.org/x/crypto/bcrypt`: password hashing/verification for admin credentials.
- `github.com/redis/go-redis/v9`: Redis client for tenant-aware rate-limiting counters.

### Frontend
- `react`: UI runtime for admin dashboard.
- `react-dom`: browser rendering for React.
- `typescript`: typed frontend development.
- `vite`: fast local dev server and production build.
- `@vitejs/plugin-react`: React JSX/Fast Refresh support in Vite.
- `@types/react`, `@types/react-dom`: TypeScript type definitions.
- `@types/node`: Node.js type definitions required by Vite/TypeScript config.

### Infrastructure
- `postgres:16-alpine`: relational storage baseline for tenant and auth data.
- `redis:7-alpine`: in-memory store for caching/rate-limiting data.

## Proxy configuration
- `PROXY_TIMEOUT_SECONDS`: upstream timeout for proxied requests (default `10`).
- `PROXY_UPSTREAMS`: tenant-safe route map loaded at startup.
  - Format: `<tenant_id>:<service>=<base_url>,<tenant_id>:<service>=<base_url>`
  - Example: `1:billing=http://localhost:18081,1:catalog=http://localhost:18082,2:billing=http://localhost:28081`

## Structured request logs
- Gateway logs are JSON (`log/slog`) and emitted once per request.
- Required fields: `tenant_id`, `route`, `status`, `latency_ms`, `request_id`.
- `X-Request-ID` is returned to clients and forwarded to upstream services.
