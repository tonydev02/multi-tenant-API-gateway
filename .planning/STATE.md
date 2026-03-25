# STATE

## Current status (2026-03-19)
- Phase 01 scaffold is implemented for backend, frontend, and local infra config.
- Backend verification (`go test`, `go vet`, `go build`) is passing.
- Frontend verification (`npm run build`) is passing.
- Docker Compose configuration and runtime startup verification are passing.
- Phase 02 implementation is in place for backend auth/tenancy and frontend login shell.
- Phase 03 tenant-aware Redis-backed rate limiting is implemented and verified.
- Phase 04 tenant-safe proxying and structured request logging are implemented and verified.
- Phase 05 admin dashboard, tenant/key management, and tenant-scoped traffic visibility are implemented and verified.
- Phase 06 repository implementation is complete:
  - backend Dockerization,
  - readiness endpoint (`/readyz`),
  - deployment templates and runbooks,
  - CI + deployment trigger workflows.

## Maintenance updates (2026-03-25)
- Route normalization used by request logging and tenant-aware rate limiting now also normalizes UUID and ULID path segments (in addition to numeric IDs) to `:id`.
- Added table-driven tests for normalized route behavior to protect key cardinality assumptions.

## Assumptions
- MVP prioritizes architecture clarity over feature depth.
- Provider-native observability is sufficient for phase 06 (OpenTelemetry deferred).
- Initial internet launch uses provider-generated domains before custom domains.
- Free-tier provider terms may change; deployment docs should be revalidated at rollout time.

## Open decisions
- Whether to keep deploy-hook based frontend automation or move to direct artifact uploads via Cloudflare API.
- When to introduce persistent analytics storage beyond in-process counters.

## Known blockers
- Live cloud credentials and staging service provisioning are required to complete internet UAT.

## Next actions
1. Provision staging resources (Render, Neon, Upstash, Cloudflare Pages).
2. Configure GitHub secrets for deploy workflows.
3. Run `deployments/scripts/smoke_public.sh` against staging URL.
4. Record live deployment evidence in Phase 06 UAT and summary docs.
