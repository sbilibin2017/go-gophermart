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
	router := chi.NewRouter()
	router.Use(
		middlewares.LoggingMiddleware,
		middlewares.GzipMiddleware,
	)
	router.Route("/api", func(r chi.Router) {
		r.Get("/orders/{number}", accrualHandler)
		r.Post("/orders", registerOrderHandler)
		r.Post("/goods", registerGoodsHandler)
	})
	return router
}
