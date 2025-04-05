package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Config interface {
	GetRunAddress() string
}

type ServerWuthRouter struct {
	*http.Server
	rtr *chi.Mux
}

func NewServerWithRouter(c Config) *ServerWuthRouter {
	rtr := chi.NewRouter()
	return &ServerWuthRouter{
		Server: &http.Server{
			Addr:    c.GetRunAddress(),
			Handler: rtr,
		},
		rtr: rtr,
	}
}

func (s *ServerWuthRouter) AddRouter(h *chi.Mux) {
	s.rtr.Mount("/", h)
}
