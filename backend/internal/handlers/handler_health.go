package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func HandelHealth(w http.ResponseWriter, r *http.Request) {
	slog.Info(
		"request received",
		"method", r.Method,
		"path", r.URL.Path,
	)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"status":  "ok",
		"service": "api-gateway",
	}

	json.NewEncoder(w).Encode(response)
}
