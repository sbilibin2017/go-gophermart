package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisteruserRegisterRoute(
	router *chi.Mux,
	prefix string,
	handler http.HandlerFunc,
	loggingMiddleware func(next http.Handler) http.Handler,
	gzipMiddleware func(next http.Handler) http.Handler,
) {
	_router := chi.NewRouter()
	_router.Use(
		loggingMiddleware,
		gzipMiddleware,
	)
	_router.Route(prefix, func(r chi.Router) {
		r.Post("/register", handler)
	})
	router.Mount("/", _router)
}
