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
    slog.ErrorContext(
        r.Context(),
        msg,
        "error", err,
    )
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(statusCode)

    errorRep := errorResponse{
        Error: msg + err.Error(),
    }

    _ = json.NewEncoder(w).Encode(errorRep)
}
