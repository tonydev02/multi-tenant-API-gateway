# PHASE-UAT: 06 Deploy and Observability

## UAT checklist

### Runtime hardening
- [x] Backend startup fails fast for missing/weak critical config (JWT secret, env safety checks).
- [x] `GET /readyz` returns `200` when dependency checks pass.
- [x] `GET /readyz` returns `503` when dependency checks fail.

### Deployment assets
- [x] Dockerfile exists for backend production deployment.
- [x] Staging/prod-lite env templates exist for backend and frontend.
- [x] Render service blueprint template exists.
- [x] Public smoke script exists and is executable.

### Observability and operations
- [x] Runbook documents required log fields and incident triage by `request_id`.
- [x] Runbook includes alert threshold guidance for free-tier operations.
- [x] Rollback runbook includes backend/frontend rollback steps and secret rotation guidance.

### CI/CD
- [x] CI workflow validates backend test/vet/build and frontend build.
- [x] Backend deploy workflow triggers Render deploy hook on main after successful CI.
- [x] Frontend deploy workflow triggers Cloudflare Pages deploy hook on main frontend changes.

### Verification
- [ ] Manual internet deployment executed on live staging URLs.
- [ ] Public smoke test executed against deployed staging service.

## Exit criteria
Phase 06 is complete when repository deployment assets, readiness checks, and CI/CD are in place, and the staged internet deployment UAT items are executed with recorded evidence.
