package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	*http.Server
}

func NewServer(addr string) *Server {
	return &Server{Server: &http.Server{Addr: addr}}
}

func (s *Server) SetHandler(router chi.Router) {
	s.Handler = router
}
