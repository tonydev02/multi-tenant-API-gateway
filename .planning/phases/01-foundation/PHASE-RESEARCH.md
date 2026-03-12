# PHASE-RESEARCH: 01 Foundation

## Objective
Define the minimal, production-style foundation for a multi-tenant API gateway SaaS repository before feature work starts.

## Architecture choices

### Monorepo layout
- Use a single repo with top-level `backend/`, `frontend/`, `deployments/`, and `.planning/`.
- Rationale: keeps backend/frontend changes reviewable together and simplifies local setup.

### Backend stack (Go)
- Use Go + standard library HTTP server for the first slice (`/health`).
- Rationale: satisfies "standard library first" and avoids premature framework selection.

### Frontend stack (React + TypeScript)
- Use React + TypeScript via Vite scaffold.
- Rationale: fast local feedback, minimal config overhead, standard ecosystem for MVP admin UI.

### Data services
- Use PostgreSQL and Redis as Docker Compose services.
- Rationale: aligns with product constraints and enables tenancy/rate-limit work in later phases.

### Local orchestration
- Use Docker Compose for Postgres/Redis and optional backend/frontend service wiring.
- Rationale: one-command local environment, consistent across contributors.

### Configuration strategy
- Read backend configuration from environment variables only.
- Rationale: required by architecture rules and portable to container/local environments.

## Dependency policy for Phase 1
- Backend dependencies: none beyond Go standard library for initial health endpoint.
- Frontend dependencies: React, React DOM, TypeScript, Vite (from scaffold).
- Infra dependencies: official Docker images for PostgreSQL and Redis.
- Rule: every non-standard dependency must be listed with one-line rationale in README and/or phase docs.

## Risks
- Frontend scaffold introduces many files; risk of noisy first commit.
- Compose config drift between local and future CI if not documented clearly.
- Missing explicit conventions early can cause inconsistent structure in later phases.

## Non-goals
- Multi-tenant business logic implementation.
- Authentication/authorization flows.
- Proxy routing and rate limiting logic.
- Observability pipelines beyond basic local logs.
