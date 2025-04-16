package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HTTPMethod string

const (
	MethodGet  HTTPMethod = http.MethodGet
	MethodPost HTTPMethod = http.MethodPost
)

func RegisterHandler(
	router *chi.Mux,
	prefix string,
	method HTTPMethod,
	handler http.HandlerFunc,
	middlewares []func(next http.Handler) http.Handler,
) {
	r := chi.NewRouter()
	r.Use(middlewares...)
	switch method {
	case MethodGet:
		r.Get(prefix, handler)
	case MethodPost:
		r.Post(prefix, handler)
	}
	router.Mount("/", r)
}
