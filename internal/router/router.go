package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HttpMethod string

const (
	MethodGet  HttpMethod = http.MethodGet
	MethodPost HttpMethod = http.MethodPost
)

func RegisterHandler(
	router chi.Router,
	prefix string,
	method HttpMethod,
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
