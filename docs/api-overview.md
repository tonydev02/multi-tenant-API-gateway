# API Overview (MVP Progress)

## Health
- `GET /health`

## Admin auth
- `POST /api/admin/login`
- `GET /api/admin/me` (JWT required)

## Tenant management
- `POST /api/admin/tenants/register`
- `GET /api/admin/tenants/current` (JWT required)
- `PATCH /api/admin/tenants/current` (JWT required)
- `DELETE /api/admin/tenants/current` (JWT required)

## API key management
- `POST /api/admin/api-keys` (JWT required)
- `GET /api/admin/api-keys` (JWT required)
- `POST /api/admin/api-keys/{id}/revoke` (JWT required)

## Admin traffic visibility
- `GET /api/admin/traffic/summary` (JWT required)

## Consumer identity
- `GET /api/consumer/whoami` (`X-API-Key` required)
- `ANY /api/consumer/proxy/{service}/{path...}` (`X-API-Key` required)

## Notes
- Admin endpoints infer tenant context from JWT claims.
- Consumer endpoints infer tenant context from API key lookup.
- Traffic summary endpoint is tenant-scoped via JWT context and does not accept client tenant IDs.
- Do not send tenant IDs from clients to choose tenancy.
- `X-Request-ID` is returned on responses and forwarded for proxy requests.
