# Security Model

## Purpose
This document defines trust boundaries and tenant isolation rules for the gateway MVP.

## Assets to protect
- Tenant data and tenant-scoped upstream access.
- Admin authentication credentials and JWT signing secret.
- Consumer API keys and key metadata.
- Request logs and operational telemetry integrity.

## Primary trust boundaries
1. Client-to-gateway boundary:
- All client-provided tenant identifiers are untrusted.
- Gateway derives tenant identity from verified credentials only.

2. Gateway-to-data boundary:
- Tenant identity for data access comes from context populated by auth middleware.
- Handlers must not accept tenant ID query/path/body inputs to choose data scope.

3. Gateway-to-upstream boundary:
- Upstream routing is resolved from `(tenant_id, service)` on the server.
- Client-provided proxy routing hints are stripped before proxying.

## Authentication model
### Admin authentication
- Mechanism: JWT bearer token in `Authorization` header.
- Validation: signature + claims via `JWT_SECRET` and configured issuer.
- Tenant source of truth: JWT claims (`tenant_id`) stored in request context.

### Consumer authentication
- Mechanism: API key in `X-API-Key`.
- Validation: key lookup through backend auth store.
- Tenant source of truth: tenant ID from key record, not from client headers.

## Tenant isolation guarantees
- Tenant context is required for protected routes.
- Rate limiting is tenant-aware.
- Proxy resolution is tenant + service scoped.
- Traffic metrics and admin data are tenant-scoped.

## Header hardening
The proxy authorization middleware strips untrusted headers before routing:
- `X-Tenant-ID`
- `X-Upstream-ID`

Gateway-controlled headers used internally:
- `X-Gateway-Tenant-ID`
- `X-Gateway-Upstream-Host`
- `X-Request-ID`

## Abuse controls
- Fixed-window tenant-aware rate limiting backed by Redis.
- Route normalization (`:id` replacement) avoids key-cardinality blowups for dynamic IDs.

## Configuration safeguards
- `JWT_SECRET` is mandatory and must be at least 32 characters.
- `ENVIRONMENT` is restricted to known values.
- `BOOTSTRAP_ON_START=true` is disallowed outside development.

## Failure behavior
- Invalid/missing auth returns `401`.
- Missing tenant context returns `401`.
- Unknown upstream mapping returns `404`.
- Upstream connectivity failures return `502`.
- Rate-limit backend issues return `503` (or `500` for invalid internal policy).

## Residual risks (MVP)
- Traffic metrics are in-memory and reset on process restart.
- No WAF/bot management layer is included in MVP core.
- No centralized SIEM pipeline is included by default.

## Recommended next hardening steps
1. Add JWT key rotation playbook and secret versioning.
2. Add persistent metrics/log shipping and retention policy.
3. Add anomaly alerts for tenant-level auth failures and 429 spikes.
4. Add optional per-tenant rate-limit policy overrides with guardrails.
