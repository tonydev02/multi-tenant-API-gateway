# Phase 06 Deployment Guide (Render + Neon + Upstash + Cloudflare Pages)

This folder contains a learning-first path to deploy the gateway to the public internet with free-tier-friendly providers.

## Environments
- `staging`: first public deployment for smoke tests and iteration.
- `prod-lite`: same topology with stricter security/runtime settings.

## Files
- `env/backend.staging.env.example`
- `env/backend.prod-lite.env.example`
- `env/frontend.staging.env.example`
- `env/frontend.prod-lite.env.example`
- `render/render.yaml`
- `scripts/smoke_public.sh`
- `runbooks/operations.md`
- `runbooks/rollback.md`

## Manual-first rollout flow
1. Copy env templates and fill real values (never commit secrets).
2. Provision Neon database and confirm `DATABASE_URL` works with `sslmode=require`.
3. Provision Upstash Redis and confirm host/password details (`REDIS_TLS=true`).
4. Deploy backend on Render using `backend/Dockerfile` and backend env variables.
5. Deploy frontend on Cloudflare Pages with `VITE_API_BASE_URL` pointing to backend URL.
6. Run `scripts/smoke_public.sh` against public endpoints.
7. After manual validation, enable GitHub Actions deploy automation.

## Required secrets for CI/CD
- `RENDER_DEPLOY_HOOK_URL`
- `CLOUDFLARE_PAGES_DEPLOY_HOOK_URL`

## Notes
- The backend runs DB migrations on startup.
- Keep `BOOTSTRAP_ON_START=false` for internet-facing staging/prod-lite.
- Rotate `JWT_SECRET` when promoting from staging to prod-lite.
