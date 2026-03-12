# PHASE-UAT: 01 Foundation

## UAT checklist

### Repo structure
- [ ] Backend, frontend, planning, and infrastructure folders exist and are discoverable.
- [ ] File naming and layout match phase plan.

### Backend
- [ ] Backend runs locally with env-based configuration.
- [ ] `GET /health` returns HTTP 200 and a simple status payload.
- [ ] Backend build/test/vet commands execute successfully.

### Frontend
- [ ] Frontend starts locally and renders placeholder dashboard.
- [ ] Frontend production build succeeds.

### Infrastructure
- [ ] `docker compose config` validates.
- [ ] PostgreSQL and Redis services start via Compose.

### Documentation
- [ ] README includes quickstart and command references.
- [ ] Added dependencies are explained with rationale.
- [ ] Planning docs reflect current status and next steps.

## Exit criteria
Phase 01 is complete when every checklist item is checked and verification commands pass in a clean local environment.
