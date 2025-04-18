package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

const gophermarGroup = "/api/user"

func RegisterUserRegisterRoute(
	r *chi.Mux,
	handler http.HandlerFunc,
	middlewares []func(http.Handler) http.Handler,
) {
	sub := chi.NewRouter()
	sub.Use(middlewares...)
	sub.Post(gophermarGroup+"/register", handler)
	r.Mount("/", sub)
}

func RegisterUserLoginRoute(
	r *chi.Mux,
	handler http.HandlerFunc,
	middlewares []func(http.Handler) http.Handler,
) {
	sub := chi.NewRouter()
	sub.Use(middlewares...)
	sub.Post(gophermarGroup+"/login", handler)
	r.Mount("/", sub)
}

func RegisterUserOrdersUploadRoute(
	r *chi.Mux,
	handler http.HandlerFunc,
	middlewares []func(http.Handler) http.Handler,
) {
	sub := chi.NewRouter()
	sub.Use(middlewares...)
	sub.Post(gophermarGroup+"/orders", handler)
	r.Mount("/", sub)
}

func RegisterUserOrderListRoute(
	r *chi.Mux,
	handler http.HandlerFunc,
	middlewares []func(http.Handler) http.Handler,
) {
	sub := chi.NewRouter()
	sub.Use(middlewares...)
	sub.Get(gophermarGroup+"/orders", handler)
	r.Mount("/", sub)
}

func RegisterUserBalanceRoute(
	r *chi.Mux,
	handler http.HandlerFunc,
	middlewares []func(http.Handler) http.Handler,
) {
	sub := chi.NewRouter()
	sub.Use(middlewares...)
	sub.Get(gophermarGroup+"/balance", handler)
	r.Mount("/", sub)
}

func RegisterUserBalanceWithdrawRoute(
	r *chi.Mux,
	handler http.HandlerFunc,
	middlewares []func(http.Handler) http.Handler,
) {
	sub := chi.NewRouter()
	sub.Use(middlewares...)
	sub.Post(gophermarGroup+"/balance/withdraw", handler)
	r.Mount("/", sub)
}

func RegisterUserWithdrawalsRoute(
	r *chi.Mux,
	handler http.HandlerFunc,
	middlewares []func(http.Handler) http.Handler,
) {
	sub := chi.NewRouter()
	sub.Use(middlewares...)
	sub.Get(gophermarGroup+"/withdrawals", handler)
	r.Mount("/", sub)
}
