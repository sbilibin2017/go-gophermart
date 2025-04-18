package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

const accrualGroup = "/api"

func RegisterGetOrderByNumberRoute(
	r *chi.Mux,
	handler http.HandlerFunc,
	middlewares []func(http.Handler) http.Handler,
) {
	sub := chi.NewRouter()
	sub.Use(middlewares...)
	sub.Get(accrualGroup+"/orders/{number}", handler)
	r.Mount("/", sub)
}

func RegisterOrdersRoute(
	r *chi.Mux,
	handler http.HandlerFunc,
	middlewares []func(http.Handler) http.Handler,
) {
	sub := chi.NewRouter()
	sub.Use(middlewares...)
	sub.Post(accrualGroup+"/orders", handler)
	r.Mount("/", sub)
}

func RegisterGoodsRoute(
	r *chi.Mux,
	handler http.HandlerFunc,
	middlewares []func(http.Handler) http.Handler,
) {
	sub := chi.NewRouter()
	sub.Use(middlewares...)
	sub.Post(accrualGroup+"/goods", handler)
	r.Mount("/", sub)
}
