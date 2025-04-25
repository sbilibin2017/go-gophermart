package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterGophermartRoutes(
	r *chi.Mux,
	prefix string,
	userRegisterHandler http.HandlerFunc,
	userLoginHandler http.HandlerFunc,
	userOrderUploadHandler http.HandlerFunc,
	userOrdersHandler http.HandlerFunc,
	userBalanceHandler http.HandlerFunc,
	userBalanceWithdrawHandler http.HandlerFunc,
	userWithdrawalsHandler http.HandlerFunc,
	authMiddleware func(http.Handler) http.Handler,
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

	_r.Group(func(r chi.Router) {
		r.Post("/register", userRegisterHandler)
		r.Post("/login", userLoginHandler)
	})

	_r.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		r.Post("/orders", userOrderUploadHandler)
		r.Get("/orders", userOrdersHandler)

		r.Get("/balance", userBalanceHandler)
		r.Post("/balance/withdraw", userBalanceWithdrawHandler)

		r.Get("/withdrawals", userWithdrawalsHandler)
	})

	r.Mount(prefix, _r)
}
