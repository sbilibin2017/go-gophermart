package utils

import "net/http"

func RespondBadRequest(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func RespondConflict(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusConflict)
}

func RespondInternalServerError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
