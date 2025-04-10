package utils

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Middleware interface {
	Apply(next http.Handler) http.Handler
}

func Wrap(h http.Handler, middlewares ...Middleware) httprouter.Handle {
	return WrapHandler(ApplyMiddlewares(h, middlewares...))
}

func ApplyMiddlewares(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, m := range middlewares {
		h = m.Apply(h)
	}
	return h
}

func WrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		h.ServeHTTP(w, r)
	}
}
