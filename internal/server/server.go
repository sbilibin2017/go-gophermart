package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	*http.Server
}

func NewServer(addr string) *Server {
	return &Server{
		Server: &http.Server{
			Addr:    addr,
			Handler: chi.NewRouter(),
		},
	}
}

func (s *Server) AddRouter(router *chi.Mux) {
	s.Handler.(*chi.Mux).Mount("/", router)
}

func (s *Server) Run(ctx context.Context) {
	go func() {
		s.ListenAndServe()
	}()
	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	s.Shutdown(shutdownCtx)
}
