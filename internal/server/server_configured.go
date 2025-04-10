package server

import "github.com/go-chi/chi/v5"

type ServerConfiguredAddresser interface {
	GetRunAddress() string
}

func NewServerConfigured(a ServerConfiguredAddresser) *ServerWithRouter {
	return NewServerWithRouter(
		NewServer(a),
		chi.NewRouter(),
	)
}
