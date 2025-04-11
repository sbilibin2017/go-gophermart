package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterGophermartPingHandler(
	router *chi.Mux,
	h http.HandlerFunc,
) {
	router.Post("/ping", h)
}
