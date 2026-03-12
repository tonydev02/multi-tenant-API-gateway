# PHASE-UAT: 01 Foundation

## UAT checklist

### Repo structure
- [x] Backend, frontend, planning, and infrastructure folders exist and are discoverable.
- [x] File naming and layout match phase plan.

### Backend
- [x] Backend runs locally with env-based configuration.
- [x] `GET /health` returns HTTP 200 and a simple status payload.
- [x] Backend build/test/vet commands execute successfully.

### Frontend
- [x] Frontend starts locally and renders placeholder dashboard.
- [x] Frontend production build succeeds.

### Infrastructure
- [x] `docker compose config` validates.
- [x] PostgreSQL and Redis services start via Compose.

### Documentation
- [x] README includes quickstart and command references.
- [x] Added dependencies are explained with rationale.
- [x] Planning docs reflect current status and next steps.

## Exit criteria
Phase 01 is complete when every checklist item is checked and verification commands pass in a clean local environment.
