package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

func RegisterRegisterRewardRoute(
	r *chi.Mux,
	prefix string,
	h http.HandlerFunc,
) {
	logger.Logger.Info("Регистрация маршрута для /goods с префиксом:", prefix)

	_r := chi.NewRouter()

	_r.Route(prefix, func(r chi.Router) {
		logger.Logger.Info("Регистрация POST-запроса на маршрут /goods")
		r.Post("/goods", h)
	})

	r.Mount("/", _r)

	logger.Logger.Info("Маршрут /goods с префиксом ", prefix, " успешно зарегистрирован.")
}
