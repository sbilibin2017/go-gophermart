package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
)

func RegisterGophermartRouter(
	router *chi.Mux,
	db *sqlx.DB,
	config *configs.JWTConfig,
	prefix string,
	userRegisterHandler http.Handler,
	userLoginHandler http.Handler,
	orderUploadHandler http.Handler,
	orderListHandler http.Handler,
	getBalanceHandler http.Handler,
	withdrawBalanceHandler http.Handler,
	getWithdrawalsHandler http.Handler,
	loggingMiddleware func(next http.Handler) http.Handler,
	gzipMiddleware func(next http.Handler) http.Handler,
	txMiddleware func(db *sqlx.DB) func(http.Handler) http.Handler,
	authMiddleware func(config *configs.JWTConfig) func(http.Handler) http.Handler,
) {
	r := chi.NewRouter()
	r.Use(
		loggingMiddleware,
		gzipMiddleware,
		txMiddleware(db),
	)

	r.Post("/register", toHandlerFunc(userRegisterHandler))
	r.Post("/login", toHandlerFunc(userLoginHandler))

	r.Route("/", func(r chi.Router) {
		r.Use(authMiddleware(config))
		r.Post("/orders/{number}", toHandlerFunc(orderUploadHandler))
		r.Get("/orders", toHandlerFunc(orderListHandler))
		r.Get("/balance", toHandlerFunc(getBalanceHandler))
		r.Post("/balance/withdraw", toHandlerFunc(withdrawBalanceHandler))
		r.Get("/withdrawals", toHandlerFunc(getWithdrawalsHandler))
	})

	router.Mount(prefix, r)
}
