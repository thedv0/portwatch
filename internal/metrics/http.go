package metrics

import (
	"encoding/json"
	"net/http"
)

// Handler returns an http.Handler that serves a JSON snapshot of r.
func Handler(r *Registry) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		snap := r.Snapshot()
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(snap); err != nil {
			http.Error(w, "encoding error", http.StatusInternalServerError)
		}
	})
}

// ListenAndServe starts a metrics HTTP server on addr (e.g. ":9090").
// It blocks until the server exits and returns any error.
func ListenAndServe(addr string, r *Registry) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", Handler(r))
	return http.ListenAndServe(addr, mux)
}
