package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
)

func RegisterAccrualRouter(
	router *chi.Mux,
	db *sqlx.DB,
	config *configs.JWTConfig,
	prefix string,
	orderInfoHandler http.Handler,
	orderRegistrationHandler http.Handler,
	goodsRegistrationHandler http.Handler,
	loggingMiddleware func(next http.Handler) http.Handler,
	gzipMiddleware func(next http.Handler) http.Handler,
	txMiddleware func(db *sqlx.DB) func(http.Handler) http.Handler,
) {
	r := chi.NewRouter()

	r.Use(
		loggingMiddleware,
		gzipMiddleware,
		txMiddleware(db),
	)

	r.Get("/orders/{number}", toHandlerFunc(orderInfoHandler))
	r.Post("/orders", toHandlerFunc(orderRegistrationHandler))
	r.Post("/goods", toHandlerFunc(goodsRegistrationHandler))

	router.Mount(prefix, r)
}
