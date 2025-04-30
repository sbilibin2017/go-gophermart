package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}

func sendTextPlainResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	w.Write([]byte(capitalize(message)))
}

func setAuthorizationHeader(w http.ResponseWriter, token string) {
	w.Header().Set("Authorization", "Bearer "+token)
}

func getURLParam(r *http.Request, param string) string {
	return chi.URLParam(r, param)
}

func decodeJSONRequest(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(v); err != nil {
		return err
	}
	return nil
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		sendTextPlainResponse(w, http.StatusInternalServerError, "Error encoding response")
	}
}

type Validator interface {
	Struct(v any) error
}

type ValidationErrorRegistry interface {
	Get(err error) *types.ValidationWithStatusCode
}

func handleValidationError(
	w http.ResponseWriter,
	req any,
	val Validator,
	valErrRegistry ValidationErrorRegistry,
) bool {
	if err := val.Struct(req); err != nil {
		valErr := valErrRegistry.Get(err)
		if valErr != nil {
			sendTextPlainResponse(w, valErr.StatusCode, valErr.Error.Error())
			return true
		}
	}
	return false
}

type HTTPErrorRegistry interface {
	Get(err error) *types.HTTPError
}

func handleServiceError(
	w http.ResponseWriter,
	err error,
	httpErrRegistry HTTPErrorRegistry,
) bool {
	if err == nil {
		return false
	}

	if httpErr := httpErrRegistry.Get(err); httpErr != nil {
		sendTextPlainResponse(w, httpErr.StatusCode, httpErr.Error.Error())
	} else {
		sendTextPlainResponse(w, http.StatusInternalServerError, "Internal server error")
	}

	return true
}

