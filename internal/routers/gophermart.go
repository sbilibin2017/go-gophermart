package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewGophermartRouter(
	authMiddleware func(http.Handler) http.HandlerFunc,
	registerHandler http.HandlerFunc,
	loginHandler http.HandlerFunc,
	uploadOrderHandler http.HandlerFunc,
	getOrdersHandler http.HandlerFunc,
	getBalanceHandler http.HandlerFunc,
	withdrawHandler http.HandlerFunc,
	getWithdrawalsHandler http.HandlerFunc,
) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/api/user/register", registerHandler)
	r.Post("/api/user/login", loginHandler)
	r.Post("/api/user/orders", authMiddleware(uploadOrderHandler))
	r.Get("/api/user/orders", authMiddleware(getOrdersHandler))
	r.Get("/api/user/balance", authMiddleware(getBalanceHandler))
	r.Post("/api/user/balance/withdraw", authMiddleware(withdrawHandler))
	r.Get("/api/user/withdrawals", authMiddleware(getWithdrawalsHandler))
	return r
}
