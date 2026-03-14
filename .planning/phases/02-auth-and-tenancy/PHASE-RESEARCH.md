# PHASE-RESEARCH: 02 Auth and Tenancy

## Objective
Introduce explicit tenancy and authentication foundations for the MVP:
- JWT auth for admin UI/API
- API keys for gateway consumers
- tenant-safe request flow that does not trust client-supplied tenant identifiers

## Architecture choices

### Tenancy model
- Add `tenants` as a first-class table and required foreign key on tenant-owned resources.
- Resolve tenant server-side from trusted credentials (JWT claims or API key lookup), not from request body/header tenant IDs.
- Rationale: enforces architecture rule that tenancy is explicit and prevents tenant spoofing.

### Admin authentication (JWT)
- Use short-lived access JWTs signed by backend secret from environment variables.
- Keep initial auth simple: username/password login endpoint for admin users.
- Rationale: meets MVP requirement while keeping implementation reviewable.

### Consumer authentication (API key)
- Store API keys as hashed values (never plaintext at rest) with key prefix for lookup.
- Support key create/revoke/list flows under tenant boundary.
- Rationale: safer key storage and clean migration path to rotation/audit in later phases.

### Backend layering
- Keep simple package boundaries:
  - transport (HTTP handlers)
  - service (business logic)
  - store (PostgreSQL/Redis access)
- Rationale: avoids unnecessary abstraction while preserving testability.

### Config and secrets
- All auth and DB settings sourced from environment variables (`JWT_SECRET`, DB connection, etc.).
- Rationale: aligns with project architecture rules and deployment portability.

## Risks
- Incorrect tenant resolution can cause cross-tenant data leakage.
- Weak JWT secret handling can compromise admin sessions.
- API key hashing/lookup design mistakes can break auth performance.
- Early schema choices may require migrations in later phases.

## Non-goals
- OAuth/OIDC or third-party identity providers.
- Refresh-token/session-device management.
- Fine-grained RBAC beyond basic admin authorization.
- External key management services (KMS/HSM) in this phase.
