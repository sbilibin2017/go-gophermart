package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
)

func NewAccrualRouter(
	accrualHandler http.HandlerFunc,
	registerOrderHandler http.HandlerFunc,
	registerGoodsHandler http.HandlerFunc,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares.LoggingMiddleware)
	r.Use(middlewares.GzipMiddleware)
	r.Route("/api", func(r chi.Router) {
		r.Get("/orders/{number}", accrualHandler)
		r.Post("/orders", registerOrderHandler)
		r.Post("/goods", registerGoodsHandler)
	})
	return r
}
