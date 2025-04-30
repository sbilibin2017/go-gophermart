package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func RegisterAccrualRouter(
	router *chi.Mux,
	db *sqlx.DB,
	prefix string,
	orderAccrualInfoHandler http.Handler,
	orderAccrualHandler http.Handler,
	goodRewardHandler http.Handler,
	loggingMiddleware func(next http.Handler) http.Handler,
	gzipMiddleware func(next http.Handler) http.Handler,
	txMiddleware func(http.Handler) http.Handler,
) {
	r := chi.NewRouter()

	r.Use(
		loggingMiddleware,
		gzipMiddleware,
		txMiddleware,
	)

	r.Get("/orders/{number}", toHandlerFunc(orderAccrualInfoHandler))
	r.Post("/orders", toHandlerFunc(orderAccrualHandler))
	r.Post("/goods", toHandlerFunc(goodRewardHandler))

	router.Mount(prefix, r)
}
