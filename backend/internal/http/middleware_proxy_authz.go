package gatewayhttp

import "net/http"

// requireProxyAuthorization strips untrusted routing hints from client input.
func requireProxyAuthorization() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Del("X-Tenant-ID")
			r.Header.Del("X-Upstream-ID")
			next.ServeHTTP(w, r)
		})
	}
}
