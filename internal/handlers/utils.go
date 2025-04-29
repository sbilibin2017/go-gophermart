package handlers

import (
	"encoding/json"
	"net/http"
)

// Decode the request body into the provided struct v
func decodeRequestBody(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(v); err != nil {
		return err
	}
	return nil
}

// Handle validation errors
func handleValidationErrorResponse(
	w http.ResponseWriter,
	err error,
) {
	validationErr := err // Assuming `handleValidationError` is not necessary; can log it directly
	if validationErr != nil {
		handleErrorResponse(w, validationErr, http.StatusBadRequest)
	}
}

// Handle errors and log the response
func handleErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	http.Error(w, capitalize(err.Error()), statusCode)
}

// Handle bad request errors and log
func handleBadRequestErrorResponse(w http.ResponseWriter) {
	http.Error(w, "Invalid request body", http.StatusBadRequest)
}

// Handle internal server errors and log
func handleInternalErrorResponse(w http.ResponseWriter) {
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

// Set the Authorization token header in the response
func setTokenHeader(w http.ResponseWriter, token string) {
	w.Header().Set("Authorization", "Bearer "+token)
}

// Send a text response with the provided message and status code
func sendTextResponse(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

// Capitalize the first letter of an error message (assuming this function exists)
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}
