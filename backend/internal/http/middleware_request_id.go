package gatewayhttp

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

const requestIDHeader = "X-Request-ID"

type requestIDContextKey struct{}

func withRequestIDContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDContextKey{}, id)
}

func requestIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDContextKey{}).(string)
	return id, ok
}

func withRequestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = generateRequestID()
			}

			w.Header().Set(requestIDHeader, requestID)
			r.Header.Set(requestIDHeader, requestID)
			next.ServeHTTP(w, r.WithContext(withRequestIDContext(r.Context(), requestID)))
		})
	}
}

func generateRequestID() string {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return "req-fallback"
	}
	return hex.EncodeToString(buf)
}
