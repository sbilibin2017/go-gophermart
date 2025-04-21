package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterHealthRoutes(
	router chi.Router,
	dbH http.HandlerFunc,
) {
	router.Get("/health", dbH)
}
