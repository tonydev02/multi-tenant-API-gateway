# PHASE-UAT: 05 Admin Dashboard

## UAT checklist

### Authenticated dashboard access
- [x] Admin can log in and navigate from auth shell to dashboard view.
- [x] Expired/invalid JWT returns user to login with clear message.

### Tenant management
- [x] Dashboard loads current tenant profile from authenticated endpoint.
- [x] Tenant name update succeeds and reflects on refresh.
- [x] Tenant operations remain scoped to JWT tenant context only.

### API key lifecycle
- [x] API key list renders with expected metadata (`id`, `name`, `prefix`, `created_at`, `revoked_at`).
- [x] Creating a key shows plaintext secret exactly once.
- [x] Revoking an active key updates list state and backend data.
- [x] Revoking a non-existent/already revoked key returns handled error state.

### Traffic/rate-limit visibility
- [x] Dashboard displays tenant-scoped traffic summary values.
- [x] Rate-limited request count increases after forcing 429 responses.
- [x] Status bucket counters change when exercising success and failure paths.
- [x] Visibility endpoint does not accept or trust client-supplied tenant IDs.

### Verification
- [x] `go test ./...` passes.
- [x] `go vet ./...` passes.
- [x] Backend build passes.
- [x] Frontend build passes.

## Exit criteria
Phase 05 is complete when dashboard tenant/key management and tenant-scoped visibility are verified through tests and runtime smoke checks.
