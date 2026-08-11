package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func registerRoutes(mux *http.ServeMux, logger *slog.Logger) {
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"ok": true}, logger)
	})
	mux.HandleFunc("GET /version", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"version": "0.1.0"}, logger)
	})
}

func writeJSON(w http.ResponseWriter, code int, body any, logger *slog.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.Error("encode response", "err", err)
	}
}
