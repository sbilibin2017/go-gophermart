package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewAccrualRouter(
	loggingMiddleware func(next http.Handler) http.Handler,
	gzipMiddleware func(next http.Handler) http.Handler,
	accrualHandler http.HandlerFunc,
	registerOrderHandler http.HandlerFunc,
	registerGoodsHandler http.HandlerFunc,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(loggingMiddleware)
	r.Use(gzipMiddleware)
	r.Route("/api", func(r chi.Router) {
		r.Get("/orders/{number}", accrualHandler)
		r.Post("/orders", registerOrderHandler)
		r.Post("/goods", registerGoodsHandler)
	})
	return r
}
