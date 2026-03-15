package gatewayhttp

import (
	"net/http"

	"github.com/namta/multi-tenant-api-gateway/backend/internal/auth"
	"github.com/namta/multi-tenant-api-gateway/backend/internal/tenant"
)

// Dependencies holds required collaborators for HTTP handlers.
type Dependencies struct {
	AuthStore      *auth.Store
	TenantStore    *tenant.Store
	JWTManager     *auth.JWTManager
	APIKeyAuth     *auth.APIKeyAuthenticator
	FrontendOrigin string
}

// NewRouter builds the HTTP router for gateway APIs.
func NewRouter(deps Dependencies) http.Handler {
	mux := http.NewServeMux()

	authMiddleware := requireAdminAuth(deps.JWTManager)
	tenantMiddleware := requireTenantContext()
	adminGuard := chainMiddleware(authMiddleware, tenantMiddleware)
	consumerGuard := requireAPIKeyAuth(deps.APIKeyAuth)

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /api/admin/login", loginHandler(deps.AuthStore, deps.JWTManager))
	mux.HandleFunc("POST /api/admin/tenants/register", registerTenantHandler(deps.TenantStore, deps.AuthStore))

	mux.Handle("GET /api/admin/me", adminGuard(meHandler()))
	mux.Handle("GET /api/admin/tenants/current", adminGuard(getCurrentTenantHandler(deps.TenantStore)))
	mux.Handle("PATCH /api/admin/tenants/current", adminGuard(updateCurrentTenantHandler(deps.TenantStore)))
	mux.Handle("DELETE /api/admin/tenants/current", adminGuard(deleteCurrentTenantHandler(deps.TenantStore)))

	mux.Handle("POST /api/admin/api-keys", adminGuard(createAPIKeyHandler(deps.AuthStore)))
	mux.Handle("GET /api/admin/api-keys", adminGuard(listAPIKeysHandler(deps.AuthStore)))
	mux.Handle("POST /api/admin/api-keys/{id}/revoke", adminGuard(revokeAPIKeyHandler(deps.AuthStore)))

	mux.Handle("GET /api/consumer/whoami", consumerGuard(consumerWhoAmIHandler(deps.TenantStore)))

	return withCORS(deps.FrontendOrigin)(mux)
}
