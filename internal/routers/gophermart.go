package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
)

func NewGophermartRouter(
	config *configs.GophermartConfig,
	loggingMiddleware func(next http.Handler) http.Handler,
	gzipMiddleware func(next http.Handler) http.Handler,
	authMiddleware func(secretKey string) func(next http.Handler) http.Handler,
	registerHandler http.HandlerFunc,
	loginHandler http.HandlerFunc,
	uploadOrderHandler http.HandlerFunc,
	getOrderHandler http.HandlerFunc,
	getBalanceHandler http.HandlerFunc,
	withdrawBalanceHandler http.HandlerFunc,
	withdrawalsHandler http.HandlerFunc,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(loggingMiddleware)
	r.Use(gzipMiddleware)
	r.Route("/api/user", func(r chi.Router) {
		r.With(authMiddleware(config.JWTSecretKey)).Post("/register", registerHandler)
		r.With(authMiddleware(config.JWTSecretKey)).Post("/login", loginHandler)
		r.With(authMiddleware(config.JWTSecretKey)).Post("/orders", uploadOrderHandler)
		r.With(authMiddleware(config.JWTSecretKey)).Get("/orders", getOrderHandler)
		r.With(authMiddleware(config.JWTSecretKey)).Get("/balance", getBalanceHandler)
		r.With(authMiddleware(config.JWTSecretKey)).Post("/balance/withdraw", withdrawBalanceHandler)
		r.With(authMiddleware(config.JWTSecretKey)).Get("/withdrawals", withdrawalsHandler)
	})
	return r
}
