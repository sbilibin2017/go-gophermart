package utils

import "net/http"

func SetJSONResponseHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func SetStatusOKResponseHeader(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}


