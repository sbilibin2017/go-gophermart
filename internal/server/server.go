package server

import (
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/log"
)

func NewServer(addr string) *http.Server {
	log.Info("Создание нового сервера", "адрес", addr)
	return &http.Server{Addr: addr}
}
