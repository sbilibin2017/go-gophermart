package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

func getURLParam(r *http.Request, name string) string {
	return chi.URLParam(r, name)
}

func decodeJSONRequest(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return err
	}
	return nil
}

func sendJSONResponse[T any](w http.ResponseWriter, resp types.APIResponse[T]) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func sendTextResponse(w http.ResponseWriter, resp *types.APIResponse[any]) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(resp.Status)
	w.Write([]byte(resp.Message))
}
