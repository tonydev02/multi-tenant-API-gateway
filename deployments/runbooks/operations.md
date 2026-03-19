# Operations Runbook

## Health and readiness
- Liveness: `GET /health`
- Dependency readiness: `GET /readyz`
- Treat `readyz != 200` as a backend incident.

## Observability checks
1. Render logs:
   - Verify structured JSON entries include `tenant_id`, `route`, `status`, `latency_ms`, `request_id`.
   - Investigate bursts of `status>=500`.
2. Render metrics:
   - Watch sustained error-rate increases and high latency.
3. Cloudflare Pages deployment logs:
   - Confirm successful builds from main branch/workflow trigger.
4. Neon and Upstash dashboards:
   - Confirm connection stability and no hard-limit throttling.

## Suggested alert thresholds (free-tier friendly)
- Backend unhealthy (`/health` or `/readyz` fails for 2 consecutive checks / 5 minutes).
- `5xx` responses > 2% for 10 minutes.
- p95 latency > 1500 ms for 10 minutes.

## Incident triage flow
1. Identify affected request(s) from user report or monitor.
2. Find `request_id` in backend logs.
3. Correlate route + tenant + status + latency.
4. Check dependency health (`/readyz`, Neon, Upstash).
5. Mitigate (rollback or env fix), then verify with smoke script.
