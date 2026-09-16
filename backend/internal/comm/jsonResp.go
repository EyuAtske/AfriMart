package comm

import (
	"encoding/json"
	"net/http"
)

func RespondwithJson(w http.ResponseWriter, r *http.Request, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		return
	}
}
