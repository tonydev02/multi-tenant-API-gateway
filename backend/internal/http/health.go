package gatewayhttp

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type healthResponse struct {
	Status string `json:"status"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{Status: "ok"})
}

func readyHandler(check func(context.Context) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if check == nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{
				"status": "degraded",
				"error":  "dependency checker is not configured",
			})
			return
		}

		checkCtx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := check(checkCtx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{
				"status": "degraded",
				"error":  "dependencies unavailable",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(healthResponse{Status: "ok"})
	}
}
