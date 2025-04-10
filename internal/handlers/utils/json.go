package utils

import (
	"encoding/json"
	"net/http"
)

func DecodeJSON[T any](w http.ResponseWriter, r *http.Request, v *T) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(v); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return err
	}
	return nil
}

func EncodeJSON[T any](w http.ResponseWriter, statusCode int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return err
	}
	w.WriteHeader(statusCode)
	return nil
}
