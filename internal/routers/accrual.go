package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterAccrualRouter(
	r *chi.Mux,
	prefix string,
	rewardRegisterHandler http.HandlerFunc,
	orderAcceptHandler http.HandlerFunc,
	orderGetByIDHandler http.HandlerFunc,
	loggingMiddleware func(http.Handler) http.Handler,
	gzipMiddleware func(http.Handler) http.Handler,
	txMiddleware func(http.Handler) http.Handler,
) {
	_r := chi.NewRouter()
	_r.Use(
		loggingMiddleware,
		gzipMiddleware,
		txMiddleware,
	)
	_r.Post("/goods", rewardRegisterHandler)
	_r.Post("/orders", orderAcceptHandler)
	_r.Get("/orders/{number}", orderGetByIDHandler)
	r.Mount(prefix, _r)
}
