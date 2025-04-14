package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterOrderRegisterRoute(
	r chi.Router,
	p string,
	h http.HandlerFunc,
	gm func(next http.Handler) http.Handler,
	lm func(next http.Handler) http.Handler,
) {
	nr := chi.NewRouter()
	r.Route(p, func(sr chi.Router) {
		sr.Use(gm, lm)
		sr.Post("/orders", h)
	})
	r.Mount("/", nr)
}
