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
	router := chi.NewRouter()
	router.Use(
		middlewares.LoggingMiddleware,
		middlewares.GzipMiddleware,
	)
	router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", registerHandler)
		r.Post("/login", loginHandler)
		r.With(middlewares.AuthMiddleware(config.JWTSecretKey)).Post("/orders", uploadOrderHandler)
		r.With(middlewares.AuthMiddleware(config.JWTSecretKey)).Get("/orders", getOrderHandler)
		r.With(middlewares.AuthMiddleware(config.JWTSecretKey)).Get("/balance", getBalanceHandler)
		r.With(middlewares.AuthMiddleware(config.JWTSecretKey)).Post("/balance/withdraw", withdrawBalanceHandler)
		r.With(middlewares.AuthMiddleware(config.JWTSecretKey)).Get("/withdrawals", withdrawalsHandler)
	})
	return router
}
