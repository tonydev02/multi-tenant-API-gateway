# Multi-Tenant API Gateway SaaS

Production-style API gateway MVP built to demonstrate why Go is a strong fit for concurrent, network-heavy backend systems.

## What Is Implemented (Through Phase 06)
- Backend: Go (`net/http`) REST service
- Frontend: React + TypeScript admin dashboard
- Data: PostgreSQL + Redis
- Local infra: Docker Compose
- Deployment assets: Dockerized backend, internet deployment runbooks, CI/CD workflows

## Key capabilities
### Admin
- Register tenant + admin user
- Login with JWT
- View/update current tenant
- Create/list/revoke API keys
- View tenant-scoped traffic summary

### Consumer
- Authenticate with `X-API-Key`
- Resolve tenant identity (`/api/consumer/whoami`)
- Proxy requests to tenant-scoped upstreams (`/api/consumer/proxy/...`)

### Ops
- Liveness endpoint: `GET /health`
- Readiness endpoint: `GET /readyz` (DB + Redis dependency checks)
- Structured JSON request logging with `request_id`

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
- `deployments/` env templates, rollout scripts, runbooks
- `.github/workflows/` CI and deploy trigger workflows
- `.planning/` phase plans, UAT, and status tracking

## Prerequisites
- Go `1.25+`
- Node.js `20+`
- npm `10+`
- Docker + Docker Compose

## Quickstart (Local)
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

## Internet Deployment (Phase 06)
Provider path: `Render + Neon + Upstash + Cloudflare Pages`.

Use:
- `deployments/README.md`
- `deployments/env/*.env.example`
- `deployments/scripts/smoke_public.sh`
- `deployments/runbooks/operations.md`
- `deployments/runbooks/rollback.md`

## Verification Commands
```bash
make backend-test
make backend-vet
make backend-build
make frontend-build
make compose-config
```

## CI/CD Workflows
- `.github/workflows/ci.yml`
- `.github/workflows/deploy-backend-render.yml`
- `.github/workflows/deploy-frontend-cloudflare.yml`

Required repository secrets:
- `RENDER_DEPLOY_HOOK_URL`
- `CLOUDFLARE_PAGES_DEPLOY_HOOK_URL`

## Config Highlights
- `ENVIRONMENT` (`development`, `staging`, `production`)
- `JWT_SECRET` (required, min 32 chars)
- `FRONTEND_ORIGIN`
- `DATABASE_URL`
- `REDIS_ADDR`, `REDIS_USERNAME`, `REDIS_PASSWORD`, `REDIS_DB`, `REDIS_TLS`
- `RATE_LIMIT_REQUESTS`, `RATE_LIMIT_WINDOW_SECONDS`
- `PROXY_TIMEOUT_SECONDS`, `PROXY_UPSTREAMS`

## Architecture and API Docs
- `docs/architecture.md`
- `docs/api-overview.md`
