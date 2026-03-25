# Observability

## Purpose
This document describes what telemetry exists today, how to interpret it, and how to debug incidents quickly.

## Logging model
The backend emits structured JSON logs per request.

### Core log fields
- `request_id`: Correlation ID for the request lifecycle.
- `tenant_id`: Tenant resolved from trusted auth context.
- `route`: Normalized route key used for observability and rate limiting.
- `status`: HTTP response status code.
- `latency_ms`: End-to-end request latency in milliseconds.
- `method`: HTTP method.
- `path`: Original incoming request path.

### Optional log fields
- `upstream_host`: Present for proxied traffic when upstream routing is resolved.

## Request ID lifecycle
1. Gateway reads incoming `X-Request-ID` if present.
2. If absent, gateway generates a request ID.
3. Response always includes `X-Request-ID`.
4. Proxy path forwards `X-Request-ID` upstream.

This enables correlation between client logs, gateway logs, and upstream logs.

## Route normalization
To avoid high-cardinality metrics/log labels, route keys are normalized:
- Numeric segment -> `:id`
- UUID segment -> `:id`
- ULID segment -> `:id`

Examples:
- `GET /api/admin/api-keys/42/revoke` -> `GET:/api/admin/api-keys/:id/revoke`
- `POST /v1/orders/123e4567-e89b-12d3-a456-426614174000/items` -> `POST:/v1/orders/:id/items`
- `DELETE /v1/sessions/01ARZ3NDEKTSV4RRFFQ69G5FAV` -> `DELETE:/v1/sessions/:id`

## Tenant traffic summary endpoint
Admin API: `GET /api/admin/traffic/summary`

### Returned aggregates
- `total_requests`
- `rate_limited_requests`
- `status_2xx`
- `status_4xx`
- `status_5xx`
- `avg_latency_ms`

### Important behavior
- Scope is tenant-only and inferred from JWT context.
- Storage is in-memory for MVP; values reset on process restart.

## Readiness and health
- `GET /health`: process liveness endpoint.
- `GET /readyz`: dependency readiness check for PostgreSQL and Redis.

## Operational checks
When investigating incidents:
1. Start from `request_id` and verify full request path.
2. Confirm `tenant_id` and normalized `route` in logs.
3. Check spikes in `429` or `5xx` by tenant.
4. Validate readiness endpoint status (`/readyz`) for dependency health.
5. Verify upstream host availability for proxy failures.

## Known observability limitations (MVP)
- No persistent metrics backend in core implementation.
- No distributed tracing exporter configured by default.
- No built-in alerting rules in application code.
