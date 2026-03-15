# PHASE-RESEARCH: 03 Rate Limiting

## Objective
Add tenant-aware request rate limiting using Redis so each tenant has isolated limits and predictable gateway behavior under load.

## Architecture choices

### Limiter strategy
- Use fixed-window counters in Redis for MVP.
- Redis key shape: `rl:{tenant_id}:{route}:{window_epoch}`.
- Rationale: simple implementation, easy to reason about, and enough for MVP control paths.

### Tenant-aware identity
- Rate limit keys must be built from server-resolved tenant ID in request context.
- Do not read tenant ID from request body/query/header.
- Rationale: preserves tenancy safety guarantees established in Phase 02.

### Scope of enforcement
- Apply limiter middleware to authenticated admin routes and consumer routes.
- Keep `/health` and login/registration outside strict limits initially.
- Rationale: avoid locking out bootstrapping/login paths while protecting core API traffic.

### Policy source
- Start with environment-driven defaults (e.g., `RATE_LIMIT_REQUESTS`, `RATE_LIMIT_WINDOW_SECONDS`).
- Allow simple route-group override in code, not database policy management yet.
- Rationale: minimal operational complexity for MVP.

### Failure behavior
- If Redis is unavailable, default to fail-closed for protected gateway routes and return 503/429-style error response.
- Rationale: avoids unlimited traffic when controls are unavailable.

## Risks
- Fixed-window burst behavior near boundary may allow short spikes.
- Redis outages can impact API availability depending on failure policy.
- Incorrect key cardinality (route granularity) can over-limit or under-limit traffic.
- Time/window math bugs can create inconsistent enforcement.

## Non-goals
- Sliding-window/log-based or token-bucket algorithms.
- Tenant self-service rate-limit policy UI.
- Global distributed coordination across multiple Redis clusters.
- Advanced adaptive or anomaly-based throttling.
