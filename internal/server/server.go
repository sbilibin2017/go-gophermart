package server

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	http.Server
}

func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		Server: http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (s *Server) Run(ctx context.Context) {
	go func() {
		s.ListenAndServe()
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Shutdown(shutdownCtx)
}
