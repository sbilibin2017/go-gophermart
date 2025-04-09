package server

import "github.com/go-chi/chi/v5"

type Addresser interface {
	GetRunAddress() string
}

func NewServerConfigured(a Addresser) *ServerWithRouter {
	return NewServerWithRouter(
		NewServer(a.GetRunAddress()),
		chi.NewRouter(),
	)
}
