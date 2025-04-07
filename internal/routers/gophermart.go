package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
)

func NewGophermartRouter(
	registerHandler http.HandlerFunc,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares.LoggingMiddleware)
	r.Use(middlewares.GzipMiddleware)
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", registerHandler)
	})
	return r
}
