package commErr

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func RespondErrorWithJson(w http.ResponseWriter, r *http.Request, statusCode int, msg string, err error) {
	if err != nil {
		slog.ErrorContext(
			r.Context(),
			msg,
			"error", err,
		)
	} else {
		slog.ErrorContext(
			r.Context(),
			msg,
		)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	errorRep := errorResponse{
		Error: msg,
	}

	_ = json.NewEncoder(w).Encode(errorRep)
}
