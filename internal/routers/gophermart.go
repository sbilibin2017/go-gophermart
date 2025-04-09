package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/jwt"
)

type JWTDecoder interface {
	Decode(tokenStr string) (*jwt.Claims, error)
}

func NewGophermartRouter(
	loggingMiddleware func(next http.Handler) http.Handler,
	gzipMiddleware func(next http.Handler) http.Handler,
	authMiddleware func(jwtDecoder JWTDecoder) func(next http.Handler) http.Handler,
	jwtDecoder JWTDecoder,
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
		r.Post("/register", registerHandler)
		r.Post("/login", loginHandler)
		r.With(authMiddleware(jwtDecoder)).Post("/orders", uploadOrderHandler)
		r.With(authMiddleware(jwtDecoder)).Get("/orders", getOrderHandler)
		r.With(authMiddleware(jwtDecoder)).Get("/balance", getBalanceHandler)
		r.With(authMiddleware(jwtDecoder)).Post("/balance/withdraw", withdrawBalanceHandler)
		r.With(authMiddleware(jwtDecoder)).Get("/withdrawals", withdrawalsHandler)
	})
	return r
}
