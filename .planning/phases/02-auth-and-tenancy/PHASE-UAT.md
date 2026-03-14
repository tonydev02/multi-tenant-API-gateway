# PHASE-UAT: 02 Auth and Tenancy

## UAT checklist

### Tenancy
- [x] Tenant model exists and is persisted in PostgreSQL.
- [x] Tenant ownership is explicit for tenant-scoped resources.
- [x] Tenant identity is server-resolved from trusted credentials.

### Admin auth (JWT)
- [x] Admin login endpoint returns JWT on valid credentials.
- [x] Invalid credentials are rejected with clear error.
- [x] Protected routes reject missing/invalid JWT.

### Consumer auth (API keys)
- [x] API keys can be created and revoked.
- [x] API keys are stored hashed (no plaintext storage).
- [x] API key auth resolves tenant correctly.

### Frontend
- [x] Basic login flow works against backend REST API.
- [x] Authenticated requests send JWT correctly.

### Verification
- [x] `go test ./...` passes.
- [x] `go vet ./...` passes.
- [x] Backend build passes.
- [x] Frontend build passes.

### Environment gap
- [ ] End-to-end smoke tests are still blocked by local host PostgreSQL role mismatch (`role "gateway" does not exist`) when running backend against `localhost:5432`.

## Exit criteria
Phase 02 is complete when all checklist items are checked and no cross-tenant access path is observed in tested flows.
