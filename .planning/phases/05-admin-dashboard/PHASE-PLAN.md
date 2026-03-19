# PHASE-PLAN: 05 Admin Dashboard

## Scope
Implement production-style admin dashboard capabilities for tenant management, API key lifecycle operations, and basic traffic/rate-limit visibility.

## Suggested subphases (optional)
- **05A Dashboard foundation**: authenticated app shell, route/section layout, shared API client auth wiring.
- **05B Tenant and key management**: tenant profile card, API key list/create/revoke interactions.
- **05C Basic visibility**: tenant traffic/rate-limit summary endpoint and dashboard widgets.

## Files to create/edit

### Frontend: dashboard UX
- `frontend/src/App.tsx` (switch from auth-only shell to authenticated dashboard shell)
- `frontend/src/features/auth/AuthShell.tsx` (retain auth entry, emit session/token to app shell)
- `frontend/src/features/dashboard/DashboardShell.tsx` (new main authenticated layout)
- `frontend/src/features/dashboard/TenantPanel.tsx` (view/update current tenant)
- `frontend/src/features/dashboard/ApiKeysPanel.tsx` (list/create/revoke key management UI)
- `frontend/src/features/dashboard/TrafficPanel.tsx` (basic traffic/rate-limit summary widgets)
- `frontend/src/lib/api.ts` (typed functions for tenant CRUD, API key ops, traffic summary)

### Backend: admin dashboard data APIs
- `backend/internal/http/router.go` (register admin visibility endpoint)
- `backend/internal/http/admin_metrics_handlers.go` (new traffic/rate-limit summary handler)
- `backend/internal/http/middleware_logging.go` (emit minimal metric events or wire collector hook)
- `backend/internal/metrics/model.go` (summary response model structs)
- `backend/internal/metrics/service.go` (in-process tenant-scoped counters for MVP visibility)
- `backend/cmd/server/main.go` (wire metrics service into HTTP dependencies)
- `backend/internal/http/response.go` (error/response helpers if needed for summary payload)

### Tests
- `frontend/src/features/dashboard/*.test.tsx` (component behavior tests where present in project setup)
- `backend/internal/metrics/service_test.go` (table-driven aggregation/reset behavior)
- `backend/internal/http/admin_metrics_handlers_test.go` (auth + tenant-scoped response checks)
- `backend/internal/http/middleware_logging_test.go` (collector integration assertions)

### Documentation
- `README.md` (dashboard capabilities and local usage flow)
- `docs/architecture.md` (admin dashboard + metrics collection flow)
- `docs/api-overview.md` (new admin traffic summary endpoint)
- `.planning/STATE.md`
- `.planning/phases/05-admin-dashboard/PHASE-UAT.md`
- `.planning/phases/05-admin-dashboard/PHASE-SUMMARY.md`

## Acceptance criteria
- Admin can log in and reach an authenticated dashboard view.
- Dashboard shows current tenant details and supports tenant name updates.
- Dashboard lists tenant API keys and supports create/revoke operations with clear success/error messaging.
- New API key plaintext is shown once in UI and not retrievable after refresh/list reload.
- Dashboard displays basic tenant traffic/rate-limit visibility from backend summary data.
- Visibility endpoint is tenant-scoped from JWT tenant context and does not accept client tenant identifiers.
- Backend and frontend builds/tests pass with documentation updated for new flows.

## Verification commands
- `cd backend && go test ./...`
- `cd backend && go vet ./...`
- `cd backend && go build ./...`
- `cd frontend && npm run build`
- `docker compose config`
- Runtime smoke checks:
  - Login from UI succeeds and dashboard sections render.
  - Tenant name update in UI is persisted and visible on reload.
  - API key create/revoke from UI reflects expected backend state changes.
  - Traffic/rate-limit widget values change after exercising protected endpoints.

## Risks
- Session/token handling mistakes can break authenticated UX or leak credentials.
- In-process traffic counters reset on restart and are not suitable for long-term analytics.
- Dashboard request fan-out can create noisy error states without careful loading/error handling.
- Visibility data can drift from expectations if middleware instrumentation misses routes.

## Non-goals
- Multi-user admin management and role-based access controls.
- Historical analytics storage, charting pipelines, or BI exports.
- API key rotation policies beyond create/list/revoke MVP flows.
- Full design-system/theming overhaul.

## UI polish addendum (2026-03-19)
1. Introduce a shared frontend stylesheet with consistent typography, spacing, color tokens, and responsive layout behavior.
2. Replace one-off inline styles in auth and dashboard components with semantic class-based styling.
3. Improve information hierarchy and visual feedback states (primary/secondary actions, status messages, table readability) without changing API behavior.
4. Re-run verification (`go test`, `go vet`, backend build, frontend build`) to keep repository health checks green.
