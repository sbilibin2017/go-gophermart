package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewAccrualRouter(
	getOrderHandler http.HandlerFunc,
	createOrderHandler http.HandlerFunc,
	registerGoodsHandler http.HandlerFunc,
) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Get("/orders/{number}", getOrderHandler)
		r.Post("/orders", createOrderHandler)
		r.Post("/goods", registerGoodsHandler)
	})
	return r
}
