package utils

import "net/http"

func ErrorInternalServerResponse(w http.ResponseWriter, err error) {
	errorResponse(w, err, http.StatusInternalServerError)
}

func ErrorBadRequestResponse(w http.ResponseWriter, err error) {
	errorResponse(w, err, http.StatusBadRequest)
}

func ErrorConflictResponse(w http.ResponseWriter, err error) {
	errorResponse(w, err, http.StatusConflict)
}

func errorResponse(w http.ResponseWriter, err error, status int) {
	http.Error(w, capitalize(buildValidationError(err).Error()), status)
}
