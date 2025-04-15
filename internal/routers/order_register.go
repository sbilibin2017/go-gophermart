package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

func RegisterOrderRegisterRoute(
	router chi.Router,
	prefix string,
	handler http.HandlerFunc,
	gzipMiddleware func(next http.Handler) http.Handler,
	loggingMiddleware func(next http.Handler) http.Handler,
) {
	logger.Logger.Infow("initializing order registration route", "prefix", prefix)

	r := chi.NewRouter()

	logger.Logger.Debug("applying gzip middleware")
	r.Use(gzipMiddleware)

	logger.Logger.Debug("applying logging middleware")
	r.Use(loggingMiddleware)

	r.Route(prefix, func(sr chi.Router) {
		sr.Post("/orders", handler)
		logger.Logger.Infow("POST /orders route registered", "path: ", prefix+"/orders")
	})

	router.Mount("/", r)

	logger.Logger.Infow("order registration route successfully mounted")
}
