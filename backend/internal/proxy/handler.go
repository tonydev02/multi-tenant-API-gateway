package proxy

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/namta/multi-tenant-api-gateway/backend/internal/auth"
)

const (
	ProxyPrefix             = "/api/consumer/proxy/"
	RequestIDHeader         = "X-Request-ID"
	ProxyUpstreamHostHeader = "X-Gateway-Upstream-Host"
)

// Handler proxies consumer traffic to tenant-safe upstream targets.
type Handler struct {
	resolver Resolver
	timeout  time.Duration
	proxy    *httputil.ReverseProxy
}

// NewHandler constructs a proxy handler with request timeout.
func NewHandler(resolver Resolver, timeout time.Duration) http.Handler {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	h := &Handler{resolver: resolver, timeout: timeout}
	h.proxy = &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			target, ok := targetFromContext(pr.In.Context())
			if !ok {
				return
			}

			_, rest := ParseProxyPath(pr.In.URL.Path)
			pr.SetURL(&target)
			pr.Out.URL.Path = normalizeTargetPath(target.Path, rest)
			pr.Out.URL.RawPath = ""
			pr.Out.Header.Del(ProxyUpstreamHostHeader)
			pr.Out.Header.Del("X-Gateway-Tenant-ID")

			if requestID := pr.In.Header.Get(RequestIDHeader); requestID != "" {
				pr.Out.Header.Set(RequestIDHeader, requestID)
			}
		},
		ErrorHandler: func(w http.ResponseWriter, _ *http.Request, _ error) {
			writeJSONError(w, http.StatusBadGateway, "upstream unavailable")
		},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID, ok := auth.TenantIDFromContext(r.Context())
		if !ok || tenantID <= 0 {
			writeJSONError(w, http.StatusUnauthorized, "tenant context missing")
			return
		}

		service, _ := ParseProxyPath(r.URL.Path)
		if service == "" {
			writeJSONError(w, http.StatusNotFound, "proxy route not found")
			return
		}

		target, err := h.resolver.Resolve(tenantID, service)
		if err != nil {
			if errors.Is(err, ErrUpstreamNotFound) {
				writeJSONError(w, http.StatusNotFound, "upstream route not found")
				return
			}
			writeJSONError(w, http.StatusBadGateway, "upstream unavailable")
			return
		}

		reqCtx, cancel := context.WithTimeout(r.Context(), h.timeout)
		defer cancel()
		r = r.WithContext(withTarget(reqCtx, target))
		r.Header.Set(ProxyUpstreamHostHeader, target.Host)
		h.proxy.ServeHTTP(w, r)
	})
}

// ParseProxyPath parses /api/consumer/proxy/{service}/{rest...}.
func ParseProxyPath(p string) (service string, rest string) {
	if !strings.HasPrefix(p, ProxyPrefix) {
		return "", ""
	}
	tail := strings.TrimPrefix(p, ProxyPrefix)
	if tail == "" {
		return "", ""
	}
	service, rest, _ = strings.Cut(tail, "/")
	service = strings.ToLower(strings.TrimSpace(service))
	if service == "" {
		return "", ""
	}
	if rest == "" {
		return service, ""
	}
	return service, "/" + rest
}

func normalizeTargetPath(basePath, rest string) string {
	if rest == "" {
		if basePath == "" {
			return "/"
		}
		return basePath
	}
	if basePath == "" || basePath == "/" {
		return rest
	}
	return path.Clean(strings.TrimSuffix(basePath, "/") + rest)
}

type targetContextKey struct{}

func withTarget(ctx context.Context, target url.URL) context.Context {
	return context.WithValue(ctx, targetContextKey{}, target)
}

func targetFromContext(ctx context.Context) (url.URL, bool) {
	target, ok := ctx.Value(targetContextKey{}).(url.URL)
	return target, ok
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
