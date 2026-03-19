# PHASE-RESEARCH: 05 Admin Dashboard

## Objective
Expand the admin UI from auth shell to a usable operations dashboard for tenant profile, API key lifecycle, and lightweight tenant traffic/rate-limit visibility.

## Architecture choices

### Frontend structure
- Keep React architecture simple with feature-local components and hooks.
- Avoid introducing global state libraries for this phase; use top-level session state plus prop-driven composition.
- Rationale: current frontend is intentionally lightweight, and Phase 05 can remain explicit and reviewable.

### Auth/session handling
- Continue JWT-based admin auth and attach token via typed API client helpers.
- Prefer memory-first session state with optional `sessionStorage` persistence only if needed for UX continuity.
- Rationale: balances usability with lower persistence risk for admin tokens.

### Tenant and API key UX
- Reuse existing backend endpoints for tenant CRUD and API key create/list/revoke.
- Surface API key plaintext only on creation success screen and never from list payloads.
- Rationale: aligns UI with backend security model where key secrets are one-time visible.

### Basic visibility model
- Add an MVP backend summary endpoint that returns per-tenant counters for:
  - total requests (recent in-process window),
  - rate-limited request count,
  - status family counts (`2xx`, `4xx`, `5xx`),
  - recent average latency.
- Populate counters from gateway middleware hooks at request completion.
- Rationale: provides actionable operational signal without introducing new storage systems in this phase.

### Tenancy safety boundary
- Traffic/limit summary must resolve tenant from trusted JWT claims.
- Do not accept tenant IDs from request parameters/body for summary lookup.
- Rationale: preserves explicit multi-tenancy guarantees across admin workflows.

## Risks
- In-process counters are reset on backend restart, which may confuse users expecting historical continuity.
- Middleware instrumentation errors can undercount/overcount request classes.
- Token expiration handling gaps can produce poor dashboard UX if refresh/login recovery is unclear.
- Parallel UI actions (e.g., create key + list refresh) can cause racey UI state if not coordinated.

## Non-goals
- Cross-tenant super-admin views.
- Long-term analytics retention in PostgreSQL/Redis.
- Advanced charting libraries or custom visualization systems.
- Client-side routing overhaul or SSR migration.
