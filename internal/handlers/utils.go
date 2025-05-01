package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

func getLoginFromContext(w http.ResponseWriter, r *http.Request) (string, error) {
	login, err := middlewares.GetLoginFromContext(r.Context())
	if err != nil {
		sendTextPlainResponse(w, types.ErrUnauthorized.Error(), http.StatusUnauthorized)
		return "", err

	}
	return login, nil
}
func sendTextPlainResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	w.Write([]byte(capitalize(message)))
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}

func decodeBody(r *http.Request, v any) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(v)
}

func decodeRequestBody(w http.ResponseWriter, r *http.Request, req interface{}) error {
	err := decodeBody(r, req)
	if err != nil {
		sendTextPlainResponse(w, types.ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return err
	}
	return nil
}

func encodeResponseBody(w http.ResponseWriter, v any, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func getURLParam(r *http.Request, param string) string {
	return chi.URLParam(r, param)
}

func setTokenHeader(w http.ResponseWriter, tokenString string) {
	w.Header().Set("Authorization", "Bearer "+tokenString)
}
