package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewAccrualRouter(
	RewardHandler http.HandlerFunc,
	orderAcceptHandler http.HandlerFunc,
	orderGetHandler http.HandlerFunc,
	middlewares []func(http.Handler) http.Handler,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares...)
	r.Post("/goods", RewardHandler)
	r.Post("/orders", orderAcceptHandler)
	r.Get("/orders/{number}", orderGetHandler)
	return r
}
