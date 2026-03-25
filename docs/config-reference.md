# Configuration Reference

## Purpose
This reference lists runtime configuration for backend, frontend, and local development support.

## Backend environment variables

| Variable | Required | Default | Used by | Notes |
|---|---|---|---|---|
| `ENVIRONMENT` | no | `development` | backend | Must be `development`, `staging`, or `production`. |
| `PORT` | no | `8080` | backend | HTTP listen port. Must be `1..65535`. |
| `FRONTEND_ORIGIN` | no | `http://localhost:5173` | backend | CORS allow-origin for admin UI. |
| `JWT_SECRET` | yes | none | backend | Minimum 32 characters; startup fails if missing/short. |
| `JWT_ISSUER` | no | `gateway-admin` | backend | JWT issuer claim. |
| `JWT_EXPIRY_MINUTES` | no | `60` | backend | Must be positive integer. |
| `DATABASE_URL` | no | `postgres://gateway:gateway@localhost:5432/gateway?sslmode=disable` | backend | Postgres DSN used for app startup, migrations, and request paths. |
| `DB_MAX_OPEN_CONNS` | no | `10` | backend | Postgres pool max open connections. Must be > 0. |
| `DB_MAX_IDLE_CONNS` | no | `5` | backend | Postgres pool max idle connections. Must be > 0. |
| `REDIS_ADDR` | no | `127.0.0.1:56379` or `REDIS_HOST:REDIS_PORT` fallback | backend | Preferred single Redis endpoint input. |
| `REDIS_HOST` | no | `127.0.0.1` | backend | Only used to compute fallback `REDIS_ADDR` when `REDIS_ADDR` is unset. |
| `REDIS_PORT` | no | `56379` | backend | Only used to compute fallback `REDIS_ADDR` when `REDIS_ADDR` is unset. |
| `REDIS_USERNAME` | no | empty | backend | Use `default` for Upstash. |
| `REDIS_PASSWORD` | no | empty | backend | Required for managed Redis in non-local environments. |
| `REDIS_DB` | no | `0` | backend | Redis logical DB index. |
| `REDIS_TLS` | no | `false` | backend | Set `true` for TLS-backed Redis endpoints. |
| `RATE_LIMIT_REQUESTS` | no | `60` | backend | Tenant-scoped fixed-window request limit. Must be > 0. |
| `RATE_LIMIT_WINDOW_SECONDS` | no | `60` | backend | Tenant-scoped fixed window in seconds. Must be > 0. |
| `PROXY_TIMEOUT_SECONDS` | no | `10` | backend | Upstream proxy timeout in seconds. Must be > 0. |
| `PROXY_UPSTREAMS` | no | empty | backend | Mapping string for tenant/service upstream targets. |
| `BOOTSTRAP_ON_START` | no | `true` | backend | Must be `false` outside development. |
| `BOOTSTRAP_TENANT_NAME` | conditional | `Acme` | backend | Required when `BOOTSTRAP_ON_START=true`. |
| `BOOTSTRAP_TENANT_SLUG` | conditional | `acme` | backend | Required when `BOOTSTRAP_ON_START=true`. |
| `BOOTSTRAP_ADMIN_EMAIL` | conditional | `admin@acme.local` | backend | Required when `BOOTSTRAP_ON_START=true`. |
| `BOOTSTRAP_ADMIN_PASSWORD` | conditional | `changeme123456` | backend | Required and minimum 12 characters when `BOOTSTRAP_ON_START=true`. |

## Frontend environment variables

| Variable | Required | Default | Used by | Notes |
|---|---|---|---|---|
| `VITE_API_BASE_URL` | no | `http://localhost:8080` | frontend | Base URL for all admin UI API calls. |

## Local development support variables
These are primarily for Docker Compose and local `.env` convenience.

| Variable | Required | Default | Used by | Notes |
|---|---|---|---|---|
| `POSTGRES_DB` | compose/local | `gateway` | docker compose | Postgres database name for local stack. |
| `POSTGRES_USER` | compose/local | `gateway` | docker compose | Postgres username for local stack. |
| `POSTGRES_PASSWORD` | compose/local | `gateway` | docker compose | Postgres password for local stack. |
| `POSTGRES_HOST` | compose/local | `localhost` | local scripts/docs | Convenience variable for local DSN composition. |
| `POSTGRES_PORT` | compose/local | `55432` | docker compose | Exposed local Postgres port. |

## `PROXY_UPSTREAMS` format
Format is a comma-separated list of mappings:
- `<tenant_id>:<service>=<base_url>`

Example:
```text
1:billing=http://localhost:18081,1:catalog=http://localhost:18082,2:billing=http://localhost:28081
```

Rules:
- `tenant_id` is an integer.
- `service` is matched from `/api/consumer/proxy/{service}/...`.
- The resolver uses `(tenant_id, service)`; client input cannot override tenant.

## Environment guidance
- `development`: bootstrap can be enabled for fast local startup.
- `staging` and `production`: set `BOOTSTRAP_ON_START=false` and provide managed Postgres/Redis credentials.
- Use unique, rotated `JWT_SECRET` values per environment.
