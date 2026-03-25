# Architecture

## Overview
This repository is a monorepo for a multi-tenant API gateway SaaS MVP:
- **Backend**: Go (`net/http`) REST API
- **Frontend**: React + TypeScript (Vite)
- **Data**: PostgreSQL + Redis
- **Local runtime**: Docker Compose
- **Internet deployment target (Phase 06)**: Render + Neon + Upstash + Cloudflare Pages

## Current implementation scope (through Phase 06)
- Health endpoint (`GET /health`)
- Dependency readiness endpoint (`GET /readyz`)
- Tenant registration and tenant CRUD (current tenant)
- Admin authentication via JWT
- Consumer authentication via API keys
- Tenant resolution from trusted credentials (JWT claims or API key lookup)
- Tenant-aware fixed-window rate limiting backed by Redis
- Tenant-safe consumer proxy routing (`/api/consumer/proxy/{service}/{path...}`)
- Request ID propagation and structured JSON request logging
- Admin dashboard for tenant profile, API key lifecycle, and traffic summary visibility
- Tenant-scoped in-memory traffic metrics (`GET /api/admin/traffic/summary`)

## Repository structure
- `backend/`
  - `cmd/server`: app entrypoint
  - `internal/config`: env-driven config loading + startup validation
  - `internal/db`: postgres connection + SQL migrations
  - `internal/http`: handlers, middleware, routing
  - `internal/metrics`: in-process tenant traffic aggregation service
- `frontend/`
  - `src/features/auth`: login/register flow
  - `src/features/dashboard`: tenant, API key, traffic summary panels
- `deployments/`
  - env templates, runbooks, Render blueprint template, public smoke script
- `.github/workflows/`
  - CI checks and deploy-trigger workflows

## Deployment design (Phase 06)
- **Backend**: containerized Go service (`backend/Dockerfile`) deployed to Render.
- **Database**: Neon PostgreSQL (`DATABASE_URL` with `sslmode=require`).
- **Rate limiting store**: Upstash Redis (`REDIS_ADDR`, `REDIS_USERNAME=default`, `REDIS_PASSWORD`, `REDIS_TLS=true`).
- **Frontend**: Cloudflare Pages with `VITE_API_BASE_URL` pointing to backend URL.

## Request flow highlights

### Readiness flow
1. Uptime checks hit `GET /health`.
2. Readiness checks hit `GET /readyz`.
3. `readyz` pings PostgreSQL and Redis and returns:
   - `200` when dependencies are reachable.
   - `503` when dependencies are unavailable.

### Logging/observability flow
1. Gateway reads or generates `X-Request-ID`.
2. Logging middleware emits one JSON log entry per request.
3. Required fields include `tenant_id`, `route`, `status`, `latency_ms`, `request_id`.
4. Route values are normalized to reduce high-cardinality keys (numeric IDs, UUIDs, and ULIDs become `:id`).
5. Operations runbooks use `request_id` for incident correlation.

## Multi-tenancy boundaries
- Tenant identity is resolved server-side from JWT/API key.
- Backend does not trust client-supplied tenant identifiers.
- Tenant-scoped operations read tenant ID from context.

## Security and config hardening
- Configuration is environment-driven.
- `JWT_SECRET` must be present and at least 32 characters.
- `ENVIRONMENT` must be `development`, `staging`, or `production`.
- `BOOTSTRAP_ON_START=true` is only allowed in `development`.
- For internet-facing deployments, use `BOOTSTRAP_ON_START=false` and rotate secrets.
