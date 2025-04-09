package server

import "github.com/go-chi/chi/v5"

func NewServerConfigured(addr string) *ServerWithRouter {
	return NewServerWithRouter(
		NewServer(addr),
		chi.NewRouter(),
	)
}
