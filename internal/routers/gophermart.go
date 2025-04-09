package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
)

func NewGophermartRouter(
	config *configs.GophermartConfig,
	registerHandler http.HandlerFunc,
	loginHandler http.HandlerFunc,
	uploadOrderHandler http.HandlerFunc,
	getOrderHandler http.HandlerFunc,
	getBalanceHandler http.HandlerFunc,
	withdrawBalanceHandler http.HandlerFunc,
	withdrawalsHandler http.HandlerFunc,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares.LoggingMiddleware)
	r.Use(middlewares.GzipMiddleware)
	r.Route("/api/user", func(r chi.Router) {
		r.With(middlewares.AuthMiddleware(config)).Post("/register", registerHandler)
		r.With(middlewares.AuthMiddleware(config)).Post("/login", loginHandler)
		r.With(middlewares.AuthMiddleware(config)).Post("/orders", uploadOrderHandler)
		r.With(middlewares.AuthMiddleware(config)).Get("/orders", getOrderHandler)
		r.With(middlewares.AuthMiddleware(config)).Get("/balance", getBalanceHandler)
		r.With(middlewares.AuthMiddleware(config)).Post("/balance/withdraw", withdrawBalanceHandler)
		r.With(middlewares.AuthMiddleware(config)).Get("/withdrawals", withdrawalsHandler)
	})
	return r
}
