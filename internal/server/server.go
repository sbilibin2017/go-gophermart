package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ServerAddresser interface {
	GetRunAddress() string
}

type Server struct {
	*http.Server
}

func NewServer(a ServerAddresser) *Server {
	return &Server{Server: &http.Server{Addr: a.GetRunAddress()}}
}

func (s *Server) SetHandler(router chi.Router) {
	s.Handler = router
}
