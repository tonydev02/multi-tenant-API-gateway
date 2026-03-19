# PHASE-RESEARCH: 06 Deploy and Observability

## Objective
Move the gateway from local-only operation to internet-accessible deployment with minimal cost and strong learning value.

## Chosen platform strategy
- **Backend hosting**: Render web service using Docker.
- **Database**: Neon serverless Postgres.
- **Rate-limit state**: Upstash Redis.
- **Frontend hosting**: Cloudflare Pages.

This path preserves current architecture (Go API + Postgres + Redis + SPA frontend) with low migration cost.

## Key decisions
- Add `GET /readyz` to separate dependency readiness from basic process liveness.
- Keep observability provider-native for this phase (Render logs/metrics + structured app logs).
- Use manual-first rollout and smoke testing before enabling deploy automation.
- Deploy automation uses provider deploy hooks to keep setup simple and auditable.

## Risks and mitigations
- **Free-tier limits can change**
  - Mitigation: keep provider abstraction in docs and fallback notes in runbooks.
- **Cold starts / idle spin-down (free services)**
  - Mitigation: separate `health` and `readyz`, add uptime check guidance.
- **Credential misconfiguration across providers**
  - Mitigation: environment templates + startup validation + smoke script.
- **Incompatible rollback after future schema changes**
  - Mitigation: rollback runbook emphasizes backward-compatible migrations.

## Non-goals
- Full OpenTelemetry pipeline.
- Multi-region failover.
- Paid SLO tooling and paging systems.
