package auth

import "context"

type contextKey string

const (
	adminClaimsKey contextKey = "adminClaims"
	tenantIDKey    contextKey = "tenantID"
)

func WithClaims(ctx context.Context, claims Claims) context.Context {
	ctx = context.WithValue(ctx, adminClaimsKey, claims)
	return context.WithValue(ctx, tenantIDKey, claims.TenantID)
}

func ClaimsFromContext(ctx context.Context) (Claims, bool) {
	claims, ok := ctx.Value(adminClaimsKey).(Claims)
	return claims, ok
}

func TenantIDFromContext(ctx context.Context) (int64, bool) {
	tenantID, ok := ctx.Value(tenantIDKey).(int64)
	return tenantID, ok
}
