# PROJECT

## Mission
Build a production-style multi-tenant API gateway SaaS in Go, with a React + TypeScript admin UI and local PostgreSQL/Redis infrastructure.

## Product goals (MVP path)
- Provide tenant-aware gateway request handling.
- Expose REST APIs for admin and management workflows.
- Support tenant-scoped rate limiting and request logging.
- Demonstrate Go strengths for concurrent network workloads.

## Technical baseline
- Backend: Go
- Frontend: React + TypeScript
- Data: PostgreSQL + Redis
- Local runtime: Docker Compose
- Auth model (planned): JWT (admin UI) + API keys (gateway consumers)

## Repository structure target
- `backend/` Go services and domain logic
- `frontend/` React admin UI
- `.planning/` project plans, phase docs, and current state
- `deployments/` deployment/runtime manifests (post-foundation)

## Engineering principles
- Small, reviewable changes.
- Standard library first for Go services.
- Explicit tenant context in data and request flows.
- Environment-variable-based configuration.
- Keep docs and verification results in sync with code.
