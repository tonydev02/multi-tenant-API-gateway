# ROADMAP

## Phase 01: Foundation
- Initialize repository scaffolding.
- Add backend service skeleton with `/health`.
- Add frontend placeholder dashboard.
- Add local PostgreSQL + Redis via Docker Compose.
- Add Makefile and README quickstart.

## Phase 02: Auth and tenancy core
- Tenant data model setup.
- API key model for consumers.
- JWT auth flow for admin UI.
- Tenant resolution middleware on backend APIs.

## Phase 03: Tenant-aware rate limiting
- Redis-backed rate limit counters.
- Tenant-scoped policy configuration.
- Enforcement middleware and tests.

## Phase 04: Proxy and request logging
- Reverse-proxy request path handling.
- Tenant-safe routing (no trust in client tenant IDs).
- Structured request logs with tenant ID, route, status, latency, request ID.

## Phase 05: Admin dashboard features
- Tenant CRUD and key management UI.
- Basic traffic/rate limit visibility.
- REST integration with backend APIs.

## Phase 06: Deploy and observability
- Containerization hardening and environment configs.
- Metrics/logging integration for operations.
- Deployment documentation and runbooks.
