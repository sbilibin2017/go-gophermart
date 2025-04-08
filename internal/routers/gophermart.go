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
	getOrdersHandler http.HandlerFunc,
	getBalanceHandler http.HandlerFunc,
	withdrawHandler http.HandlerFunc,
	getWithdrawalsHandler http.HandlerFunc,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares.LoggingMiddleware)
	r.Use(middlewares.GzipMiddleware)

	r.Route("/api/user", func(r chi.Router) {
		if registerHandler != nil {
			r.Post("/register", registerHandler)
		}
		if loginHandler != nil {
			r.Post("/login", loginHandler)
		}

		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware(config))
			if uploadOrderHandler != nil {
				r.Post("/orders", uploadOrderHandler)
			}
			if getOrdersHandler != nil {
				r.Get("/orders", getOrdersHandler)
			}
			if getBalanceHandler != nil {
				r.Get("/balance", getBalanceHandler)
			}
			if withdrawHandler != nil {
				r.Post("/balance/withdraw", withdrawHandler)
			}
			if getWithdrawalsHandler != nil {
				r.Get("/withdrawals", getWithdrawalsHandler)
			}
		})
	})
	return r
}
