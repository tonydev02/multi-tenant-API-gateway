# Rollback and Secrets Runbook

## Secrets management
- Store production secrets only in provider secret managers / GitHub Secrets.
- Never place secrets in repo-tracked files.
- Rotate at minimum:
  - `JWT_SECRET` when credentials leak is suspected.
  - Database and Redis credentials after compromise events.

## Backend rollback (Render)
1. Open Render service deploy history.
2. Roll back to last known healthy deploy.
3. Verify `GET /health` and `GET /readyz`.
4. Re-run `deployments/scripts/smoke_public.sh`.

## Frontend rollback (Cloudflare Pages)
1. Promote the previous successful deployment.
2. Confirm login/dashboard behavior with live backend.

## Database safety
- Since migrations run on startup, use backward-compatible schema changes for future phases.
- If rollback is blocked by schema drift, deploy a compatibility patch before full rollback.
