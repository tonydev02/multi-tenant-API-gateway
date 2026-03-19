# PHASE-PLAN: 06 Deploy and Observability

## Scope
Ship internet-ready deployment assets and operational practices for the MVP with a learning-first rollout path:
- Manual deployment to free-tier-friendly providers.
- CI/CD automation after manual validation.
- Basic production hardening and dependency readiness checks.

## Deployment topology
- Backend: Render (Docker deploy from `backend/Dockerfile`)
- Database: Neon Postgres
- Redis: Upstash Redis
- Frontend: Cloudflare Pages

## Implementation breakdown
- **06A Runtime hardening**
  - Add backend container image definition.
  - Enforce startup validation for environment/safety defaults.
  - Add readiness endpoint (`GET /readyz`) for DB/Redis dependency checks.
- **06B Internet rollout assets**
  - Add environment templates for `staging` and `prod-lite`.
  - Add manual rollout guide and public smoke script.
  - Add operations + rollback runbooks.
- **06C CI/CD automation**
  - Add CI workflow for backend/frontend build verification.
  - Add deploy-trigger workflow for Render backend.
  - Add deploy-trigger workflow for Cloudflare Pages frontend.

## Files created/updated
- Backend/runtime
  - `backend/Dockerfile`
  - `backend/.dockerignore`
  - `backend/internal/config/config.go`
  - `backend/internal/http/health.go`
  - `backend/internal/http/router.go`
  - `backend/cmd/server/main.go`
  - `backend/internal/http/health_test.go`
- Deployment assets
  - `deployments/README.md`
  - `deployments/env/backend.staging.env.example`
  - `deployments/env/backend.prod-lite.env.example`
  - `deployments/env/frontend.staging.env.example`
  - `deployments/env/frontend.prod-lite.env.example`
  - `deployments/render/render.yaml`
  - `deployments/scripts/smoke_public.sh`
  - `deployments/runbooks/operations.md`
  - `deployments/runbooks/rollback.md`
- CI/CD
  - `.github/workflows/ci.yml`
  - `.github/workflows/deploy-backend-render.yml`
  - `.github/workflows/deploy-frontend-cloudflare.yml`
- Docs/planning
  - `README.md`
  - `docs/architecture.md`
  - `docs/api-overview.md`
  - `.planning/STATE.md`
  - `.planning/phases/06-deploy-and-observability/PHASE-RESEARCH.md`
  - `.planning/phases/06-deploy-and-observability/PHASE-UAT.md`
  - `.planning/phases/06-deploy-and-observability/PHASE-SUMMARY.md`

## Acceptance criteria
- Backend builds as a Docker image and serves `/health` + `/readyz`.
- Readiness endpoint returns `503` when dependencies are unavailable.
- Deployment templates exist for both staging and prod-lite.
- Manual rollout runbook and smoke script can validate internet deployment end-to-end.
- CI workflow validates backend/frontend build quality on PRs and main.
- Deploy workflows can trigger backend/frontend deployment from main with configured secrets.
- Docs are updated to reflect phase 06 architecture and operations.

## Verification commands
- `cd backend && go test ./...`
- `cd backend && go vet ./...`
- `cd backend && go build ./...`
- `cd frontend && npm run build`
- `docker compose config`
