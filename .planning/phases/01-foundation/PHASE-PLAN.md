# PHASE-PLAN: 01 Foundation

## Scope
Create the initial repository scaffolding and runnable local baseline for backend, frontend, and dependencies.

## Files to create/edit

### Repository and docs
- `README.md` (project intro, run instructions, dependency explanations)
- `.planning/PROJECT.md` (project charter and structure)
- `.planning/ROADMAP.md` (phase sequencing)
- `.planning/STATE.md` (current state, assumptions, blockers)

### Backend (Go)
- `backend/go.mod`
- `backend/cmd/server/main.go`
- `backend/internal/config/config.go` (env loading, minimal)
- `backend/internal/http/health.go`
- `backend/internal/http/router.go`
- `backend/Makefile` or root `Makefile` targets for backend tasks

### Frontend (React + TypeScript)
- `frontend/package.json`
- `frontend/tsconfig.json`
- `frontend/vite.config.ts`
- `frontend/index.html`
- `frontend/src/main.tsx`
- `frontend/src/App.tsx` (placeholder dashboard)

### Local infrastructure
- `docker-compose.yml` (postgres + redis + optional app service stubs)
- `.env.example` (document required env vars)
- Root `Makefile` with setup/build/test convenience targets

## Implementation order
1. Create root docs and planning docs to lock conventions.
2. Scaffold backend minimal HTTP server with `/health`.
3. Scaffold frontend placeholder dashboard.
4. Add Docker Compose for PostgreSQL + Redis.
5. Wire Makefile targets for dev/build/test verification.
6. Verify all required commands pass (or document expected temporary failures).

## Acceptance criteria
- Repo contains backend/frontend/data scaffolding in clear directories.
- `docker compose up` starts PostgreSQL and Redis successfully.
- Backend starts and returns `200` on `GET /health`.
- Frontend starts and shows a placeholder dashboard page.
- Root README includes:
  - setup steps
  - run instructions
  - dependency list with justification for each dependency
- Root Makefile contains targets for at least: backend run/build/test/vet, frontend build, compose up/down.
- Planning docs are updated: `PROJECT.md`, `ROADMAP.md`, `STATE.md`, and this phase folder docs.

## Verification commands
- `go test ./...` (from `backend/`)
- `go vet ./...` (from `backend/`)
- `go build ./...` (from `backend/`)
- `npm run build` (from `frontend/`)
- `docker compose config`
- Optional smoke checks:
  - `curl http://localhost:<backend-port>/health`

## Risks
- Tooling versions (Go/Node/Docker) may differ by machine.
- Frontend build can fail if lockfile/node version mismatch.
- Compose service naming/ports may conflict with local running services.

## Non-goals
- Tenant resolution middleware.
- API key or JWT auth.
- Rate limiting.
- Gateway proxy routing.
- Production deployment manifests.
