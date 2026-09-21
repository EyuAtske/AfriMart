package commErr

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func RespondErrorWithJson(w http.ResponseWriter, statusCode int, err string) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(statusCode)

    errorRep := errorResponse{
        Error: err,
    }

    _ = json.NewEncoder(w).Encode(errorRep)
}
