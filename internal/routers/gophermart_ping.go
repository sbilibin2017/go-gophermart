package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterGophermartPingRoute(
	router *chi.Mux,
	h http.HandlerFunc,
) {
	router.Get("/ping", h)
}
