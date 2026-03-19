# PHASE-UAT: 05 Admin Dashboard

## UAT checklist

### Authenticated dashboard access
- [ ] Admin can log in and navigate from auth shell to dashboard view.
- [ ] Expired/invalid JWT returns user to login with clear message.

### Tenant management
- [ ] Dashboard loads current tenant profile from authenticated endpoint.
- [ ] Tenant name update succeeds and reflects on refresh.
- [ ] Tenant operations remain scoped to JWT tenant context only.

### API key lifecycle
- [ ] API key list renders with expected metadata (`id`, `name`, `prefix`, `created_at`, `revoked_at`).
- [ ] Creating a key shows plaintext secret exactly once.
- [ ] Revoking an active key updates list state and backend data.
- [ ] Revoking a non-existent/already revoked key returns handled error state.

### Traffic/rate-limit visibility
- [ ] Dashboard displays tenant-scoped traffic summary values.
- [ ] Rate-limited request count increases after forcing 429 responses.
- [ ] Status bucket counters change when exercising success and failure paths.
- [ ] Visibility endpoint does not accept or trust client-supplied tenant IDs.

### Verification
- [ ] `go test ./...` passes.
- [ ] `go vet ./...` passes.
- [ ] Backend build passes.
- [ ] Frontend build passes.

## Exit criteria
Phase 05 is complete when dashboard tenant/key management and tenant-scoped visibility are verified through tests and runtime smoke checks.
