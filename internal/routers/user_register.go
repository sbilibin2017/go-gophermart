package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
)

func RegisterUserRegisterHandler(
	router *chi.Mux,
	config *configs.GophermartConfig,
	h http.HandlerFunc,
) {
	router.With(
		middlewares.LoggingMiddleware,
		middlewares.GzipMiddleware,
	).Post("/api/user/register", h)
}
