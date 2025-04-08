package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
)

func NewAccrualRouter(
	getOrderAccrualHandler http.HandlerFunc,
	registerOrderHandler http.HandlerFunc,
	registerGoodsHandler http.HandlerFunc,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares.LoggingMiddleware)
	r.Use(middlewares.GzipMiddleware)

	r.Route("/api", func(r chi.Router) {
		r.Route("/orders", func(r chi.Router) {
			if getOrderAccrualHandler != nil {
				r.Get("/{number}", getOrderAccrualHandler)
			}
			if registerOrderHandler != nil {
				r.Post("/", registerOrderHandler)
			}
		})
		if registerGoodsHandler != nil {
			r.Post("/goods", registerGoodsHandler)
		}
	})

	return r
}
