package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
)

func NewGophermartRouter(
	config *configs.JWTConfig,
	registerHandler http.HandlerFunc,
	loginHandler http.HandlerFunc,
	uploadOrderHandler http.HandlerFunc,
	getOrdersHandler http.HandlerFunc,
	getBalanceHandler http.HandlerFunc,
	withdrawHandler http.HandlerFunc,
	getWithdrawalsHandler http.HandlerFunc,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares.LoggingMiddleware)
	r.Use(middlewares.GzipMiddleware)
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", registerHandler)
		r.Post("/login", loginHandler)
		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware(config))
			r.Post("/orders", uploadOrderHandler)
			r.Get("/orders", getOrdersHandler)
			r.Get("/balance", getBalanceHandler)
			r.Post("/balance/withdraw", withdrawHandler)
			r.Get("/withdrawals", getWithdrawalsHandler)
		})
	})
	return r
}
