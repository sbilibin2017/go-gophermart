package utils

import "net/http"

func RespondWithError(w http.ResponseWriter, err error, status int) {
	http.Error(w, err.Error(), status)
}
