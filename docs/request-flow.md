# Request Flow

## Purpose
This document explains the end-to-end request lifecycle in the gateway so contributors can reason about correctness, tenancy boundaries, and performance behavior.

## High-level pipeline
Most HTTP requests follow this sequence:
1. `withRequestID`
2. `withRequestLogging`
3. `withCORS`
4. route-specific auth and tenant middleware
5. optional tenant-aware rate limiting
6. handler logic (admin/consumer/proxy)

## Middleware ordering
The router wraps handlers in this order:
1. `withRequestID` (outermost)
2. `withRequestLogging`
3. `withCORS`
4. route middleware chain from `chainMiddleware(...)`

This means every request gets a request ID and logging, including failures from auth and rate limiting.

## Admin request flow example
Example: `GET /api/admin/traffic/summary`
1. `withRequestID` reads `X-Request-ID` or generates one.
2. `withCORS` validates browser origin policy.
3. `requireAdminAuth` validates JWT from `Authorization: Bearer ...`.
4. JWT claims are attached to context (`auth.WithClaims(...)`).
5. `requireTenantContext` ensures tenant ID exists in trusted auth claims.
6. `requireTenantRateLimit` checks tenant + normalized route against Redis-backed policy.
7. `trafficSummaryHandler` reads tenant ID from context and returns tenant-scoped metrics.
8. `withRequestLogging` emits final structured log fields.

## Consumer identity flow example
Example: `GET /api/consumer/whoami`
1. `requireAPIKeyAuth` validates `X-API-Key`.
2. API key lookup returns tenant ID from backend storage.
3. Tenant claims are attached to context server-side.
4. `requireTenantContext` ensures tenant is resolved.
5. `requireTenantRateLimit` applies consumer policy per tenant.
6. Handler returns tenant identity from trusted context.

## Consumer proxy flow example
Example: `GET /api/consumer/proxy/billing/invoices/123`
1. `requireAPIKeyAuth` authenticates the key.
2. `requireTenantContext` confirms tenant identity in context.
3. `requireTenantRateLimit` checks tenant + normalized route key.
4. `requireProxyAuthorization` strips untrusted client routing hints:
- `X-Tenant-ID`
- `X-Upstream-ID`
5. Proxy handler parses `/api/consumer/proxy/{service}/{rest...}`.
6. Resolver looks up upstream using `(tenant_id, service)`.
7. Reverse proxy rewrites target URL and path.
8. Gateway forwards `X-Request-ID` to upstream.
9. Response returns through gateway and is logged.

## Route normalization in flow
Rate limiting and request logs both use `normalizedRoute(r)`.
- Numeric segments are replaced with `:id`.
- UUID segments are replaced with `:id`.
- ULID segments are replaced with `:id`.

Example:
- raw path: `/v1/orders/123e4567-e89b-12d3-a456-426614174000/items`
- normalized route: `GET:/v1/orders/:id/items`

## Error behavior summary
- Missing/invalid auth: `401`
- Missing tenant context: `401`
- Rate limit exceeded: `429`
- Rate limiter unavailable: `503` (or `500` for invalid internal policy)
- Unknown proxy service: `404`
- Upstream unavailable: `502`

## Why this structure matters
- Tenant identity is always derived from trusted credentials.
- Rate limiting and logs remain cardinality-safe under dynamic IDs.
- Request IDs provide traceability across gateway and upstream services.
