package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterAccrualRoutes(
	router chi.Router,
	getOrderByNumberHandler http.HandlerFunc,
	registerOrderHandler http.HandlerFunc,
	registerRewardHandler http.HandlerFunc,
) {
	router.Route("/api", func(r chi.Router) {
		r.Get("/orders/{number}", getOrderByNumberHandler)
		r.Post("/orders", registerOrderHandler)
		r.Post("/goods", registerRewardHandler)
	})
}
