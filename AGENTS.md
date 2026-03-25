# AGENTS.md

## Project mission
Build a production-style multi-tenant API gateway SaaS in Go.
The project must demonstrate why Go is a strong choice for concurrent, network-heavy backend systems.

## Working style
- Prefer small, reviewable changes.
- Before implementing complex work, write or update a plan in `.planning/`.
- Keep docs in sync with code changes.
- Do not introduce unnecessary abstractions early.
- Favor simple, explicit Go over clever patterns.
- Preserve a working build at all times.
- Update documentation in `/docs` whenever codebase changes affect architecture, APIs, configuration, or developer workflow.

## Git workflow (interview-ready)
- Use `main` as protected production branch; do not push directly to `main` for normal work.
- Create short-lived branches from `main` for every change.
- Branch naming must be one of: `feat/<area>-<short-description>`, `fix/<area>-<short-description>`, `docs/<area>-<short-description>`, `chore/<area>-<short-description>`.
- Open a pull request to `main` and require CI checks to pass before merge.
- Use Conventional Commits for commit messages (for example: `feat(proxy): add tenant-safe upstream cache`).
- Keep PRs small and focused (prefer one concern per PR).
- Squash-merge PRs to keep `main` history clean and readable.
- Delete merged branches to keep the repository tidy.
- Use `hotfix/*` branches only for urgent production fixes, then merge via PR back to `main`.

## Tech constraints
- Backend: Go
- Frontend: React + TypeScript
- Data: PostgreSQL, Redis
- Local environment: Docker Compose
- API style: REST for MVP
- Auth: JWT for admin UI, API keys for gateway consumers

## Go engineering rules
- Use standard library first where reasonable.
- Keep packages small and cohesive.
- Pass context.Context through request flows.
- Return explicit errors; avoid hidden control flow.
- Add table-driven tests for business logic.
- Use interfaces only where they improve testability or boundaries.
- Avoid premature generics.

## Architecture rules
- Multi-tenancy must be explicit in data model and request flow.
- Rate limiting must be tenant-aware.
- Proxy layer must not trust client-supplied tenant identifiers.
- Request logs must include tenant ID, route, status, latency, and request ID.
- Configuration must come from environment variables.

## Safety rails
- Never delete large sections of code without explaining why.
- Never add a dependency without stating the reason.
- Never skip tests, lint, or build verification after code changes.
- If requirements are unclear, update `.planning/STATE.md` with assumptions.

## Verification
After each meaningful change:
- run `go test ./...`
- run `go vet ./...`
- run backend build
- run frontend build
- update relevant planning docs
