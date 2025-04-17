package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RegisterRegisterRewardRoute регистрирует маршрут для регистрации вознаграждения.
func RegisterRegisterRewardRoute(
	prefix string,
	r *chi.Mux,
	h http.HandlerFunc,
	txMW func(http.Handler) http.Handler,
) {
	_r := chi.NewRouter()
	_r.Use(txMW)
	_r.Post("/goods", h)
	r.Mount(prefix, _r)
}
