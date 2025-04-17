package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

func NewServer(addr string, router *chi.Mux) *http.Server {
	logger.Logger.Info("Создание нового сервера с адресом:", addr)
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	return server
}
