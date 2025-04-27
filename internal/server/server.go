package server

import "net/http"

func NewServer(addr string) *http.Server {
	return &http.Server{Addr: addr}
}
