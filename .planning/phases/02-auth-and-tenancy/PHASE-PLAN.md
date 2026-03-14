# PHASE-PLAN: 02 Auth and Tenancy

## Scope
Implement MVP authentication and tenancy core for backend APIs and foundational admin UI wiring.

## Files to create/edit

### Backend: domain and storage
- `backend/internal/tenant/model.go`
- `backend/internal/tenant/store_postgres.go`
- `backend/internal/auth/model.go`
- `backend/internal/auth/password.go`
- `backend/internal/auth/jwt.go`
- `backend/internal/auth/apikey.go`
- `backend/internal/auth/store_postgres.go`
- `backend/internal/db/postgres.go`
- `backend/internal/db/migrations/` (initial tenant/auth migration SQL)

### Backend: API and middleware
- `backend/internal/http/router.go` (register auth/tenant routes)
- `backend/internal/http/auth_handlers.go`
- `backend/internal/http/tenant_handlers.go`
- `backend/internal/http/middleware_auth.go`
- `backend/internal/http/middleware_tenant.go`

### Backend: config
- `backend/internal/config/config.go` (JWT and DB settings)
- `.env.example` (add JWT/DB env variables)

### Frontend
- `frontend/src/App.tsx` (replace placeholder with auth-aware shell)
- `frontend/src/features/auth/` (login form, token handling)
- `frontend/src/lib/api.ts` (REST client with auth header)

### Documentation
- `README.md` (auth setup and API usage)
- `.planning/STATE.md` (assumptions, progress)
- `.planning/phases/02-auth-and-tenancy/PHASE-UAT.md`
- `.planning/phases/02-auth-and-tenancy/PHASE-SUMMARY.md`

## Acceptance criteria
- Tenancy is explicit in schema and request flow.
- Admin login endpoint issues valid JWT for authenticated admin user.
- Protected admin endpoints reject missing/invalid JWT.
- API key auth path identifies tenant from key lookup (not client tenant input).
- API key plaintext is never stored in database.
- Basic tenant CRUD endpoints exist with tenant-safe behavior.
- README documents env vars and auth usage for local testing.

## Verification commands
- `cd backend && go test ./...`
- `cd backend && go vet ./...`
- `cd backend && go build ./...`
- `cd frontend && npm run build`
- `docker compose config`
- Integration smoke checks (local):
  - login -> obtain JWT -> call protected route
  - create/revoke API key -> validate tenant resolution path

## Risks
- Tenant context propagation bugs across middleware/handlers.
- Migration mismatch between local DB state and code assumptions.
- Security regressions from incorrect JWT validation.
- Scope creep into advanced auth features.

## Non-goals
- SSO, MFA, advanced RBAC.
- Full admin UX polish.
- Rate limiting implementation (Phase 03).
- Proxy routing/logging enhancements (Phase 04).
