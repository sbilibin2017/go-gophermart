package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/log"
)

func NewServer(addr string, router *chi.Mux) *http.Server {
	log.Info("Создание нового сервера", "адрес", addr)
	return &http.Server{Addr: addr, Handler: router}
}
