# Architecture

## Overview
This repository is a monorepo for a multi-tenant API gateway SaaS MVP:
- **Backend**: Go (`net/http`) REST API
- **Frontend**: React + TypeScript (Vite)
- **Data**: PostgreSQL (system of record) + Redis (reserved for rate limiting/caching phases)
- **Local runtime**: Docker Compose

## Current implementation scope (through Phase 02)
- Health endpoint (`GET /health`)
- Tenant registration and tenant CRUD (current tenant)
- Admin authentication via JWT
- Consumer authentication via API keys
- Tenant resolution from trusted credentials (JWT claims or API key lookup)

## Repository structure
- `backend/`
  - `cmd/server`: app entrypoint
  - `internal/config`: env-driven config loading
  - `internal/db`: postgres connection + SQL migrations
  - `internal/tenant`: tenant model + postgres store
  - `internal/auth`: password hashing, JWT, API key logic, auth store
  - `internal/http`: handlers, middleware, routing
- `frontend/`
  - `src/features/auth`: basic login/register shell
  - `src/lib/api.ts`: REST client helpers
- `docker-compose.yml`: postgres + redis services
- `.planning/`: phased planning and status docs

## Request flow

### Admin flow (JWT)
1. Admin logs in via `POST /api/admin/login`.
2. Backend validates credentials (`bcrypt`) and issues signed JWT.
3. Protected admin routes use `Authorization: Bearer <token>`.
4. Middleware validates token and injects claims + tenant ID into request context.

### Consumer flow (API key)
1. Admin creates API key via `POST /api/admin/api-keys`.
2. Backend stores only **hashed** key and prefix in DB; plaintext is returned once.
3. Consumer sends `X-API-Key` on protected consumer route.
4. Middleware resolves prefix, compares hash, and injects tenant context.

## Multi-tenancy boundaries
- Tenant identity is resolved server-side from JWT/API key.
- Backend does not trust client-supplied tenant identifiers.
- Tenant-scoped operations read tenant ID from context.
- Core tenant-aware tables include `tenant_id` foreign keys.

## Data model (current)
Defined in `backend/internal/db/migrations/0001_auth_tenancy.sql`:
- `tenants`
- `admin_users` (with `tenant_id` FK)
- `api_keys` (with `tenant_id` FK, hashed keys, revoke timestamp)
- `schema_migrations` (migration tracking)

## Configuration
Primary backend config is environment-based via `.env`:
- `PORT`
- `DATABASE_URL`
- `JWT_SECRET`, `JWT_ISSUER`, `JWT_EXPIRY_MINUTES`
- bootstrap admin/tenant values
- Local Docker defaults intentionally use non-standard host ports to avoid conflicts:
  - PostgreSQL: `55432`
  - Redis: `56379`

## Security choices (MVP)
- Passwords hashed with `bcrypt`
- API keys hashed with SHA-256 and verified with constant-time compare
- JWT signed with HMAC-SHA256
- Plain API key secrets are never persisted

## Next architectural milestones
- Phase 03: tenant-aware rate limiting (Redis)
- Phase 04: proxy layer + structured request logging
- Phase 05+: admin UX expansion and operational hardening
