package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/services"
)

func decodeJSONRequest(w http.ResponseWriter, r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return err
	}
	return nil
}

func encodeJSONResponse(w http.ResponseWriter, v any, status *services.APIStatus) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status.StatusCode)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func handleError(w http.ResponseWriter, errStatus *services.APIStatus) {
	if errStatus != nil {
		http.Error(w, errStatus.Message, errStatus.StatusCode)
	}
}

func sendTextPlainResponse(w http.ResponseWriter, status *services.APIStatus) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status.StatusCode)
	w.Write([]byte(status.Message))
}

func getURLParam(r *http.Request, name string) string {
	return chi.URLParam(r, name)
}

func setAuthorizationHeader(w http.ResponseWriter, token string) {
	w.Header().Set("Authorization", "Bearer "+token)
}

func getUserLoginFromContext(
	w http.ResponseWriter,
	r *http.Request,
	jwtPayloadGetter func(ctx context.Context) (map[string]any, error),
) *string {
	payload, err := jwtPayloadGetter(r.Context())
	if err != nil {
		handleError(w, &services.APIStatus{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized",
		})
		return nil
	}
	login := payload["login"].(string)
	return &login
}
