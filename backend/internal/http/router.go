package gatewayhttp

import "net/http"

// NewRouter builds the HTTP router for gateway APIs.
func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	return mux
}
