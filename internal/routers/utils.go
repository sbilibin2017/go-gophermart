package routers

import "net/http"

func toHandlerFunc(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(h.ServeHTTP)
}
