# PHASE-RESEARCH: 04 Proxy and Logging

## Objective
Introduce a tenant-safe proxy layer and structured request logging that is ready for production-style operational visibility.

## Architecture choices

### Proxy model
- Use Go `httputil.ReverseProxy` as the core proxy mechanism.
- Resolve upstream target server-side from trusted route + tenant policy data.
- Never accept client-supplied upstream host or tenant IDs for routing decisions.
- Rationale: standard library proxy is reliable and keeps implementation simple.

### Tenant safety boundary
- Proxy middleware requires authenticated tenant context from prior auth/rate-limit middleware.
- Upstream selection is based on server-side config/store keyed by tenant and route.
- Rationale: enforces architecture rule that proxy must not trust client tenant identifiers.

### Logging model
- Add structured JSON logs emitted once per request at the gateway edge.
- Required fields per request log:
  - `request_id`
  - `tenant_id`
  - `route`
  - `status`
  - `latency_ms`
- Include method/path/upstream host as supporting fields.
- Rationale: enables debugging, tracing, and SRE workflows without heavy dependencies.

### Request ID strategy
- Generate request ID when absent (`X-Request-ID`), propagate to upstream, and include in logs/response headers.
- Rationale: single correlation key across client, gateway, and upstream logs.

### Failure behavior
- Upstream resolution failure: return `404`/`502` with structured error response.
- Upstream timeout/network error: return `502` and log error context.
- Rationale: clear failure classification for operators and tests.

## Risks
- Proxy misconfiguration can accidentally route traffic to wrong upstreams.
- Missing log fields reduce debuggability and break UAT expectations.
- Large/verbose logs can increase IO overhead at high request rates.
- Error handling paths may leak internal details if not sanitized.

## Non-goals
- Full distributed tracing stack (OpenTelemetry exporter).
- Complex dynamic service discovery systems.
- Log shipping/aggregation platform integration.
- Advanced traffic shaping/canary routing.
