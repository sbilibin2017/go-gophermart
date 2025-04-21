package helpers

import (
	"encoding/json"
	"io"
	"net/http"
)

func DecodeRequestBody(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		if err == io.EOF {
			http.Error(w, "Request body is empty", http.StatusBadRequest)
		} else {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		}
		defer r.Body.Close()
		return err
	}
	defer r.Body.Close()
	return nil
}

func EncodeResponseBody(w http.ResponseWriter, v interface{}, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return err
	}
	return nil
}
