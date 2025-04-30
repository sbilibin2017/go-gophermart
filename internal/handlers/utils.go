package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}

func sendPlainTextResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func getURLParam(r *http.Request, param string) string {
	return chi.URLParam(r, param)
}
