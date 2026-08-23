package handlers

import "net/http"


func HandelProducts(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"ok"}`))
}