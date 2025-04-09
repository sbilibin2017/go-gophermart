package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/log"
)

type HTTPServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
	SetHandler(router chi.Router)
}

type ServerWithRouter struct {
	server HTTPServer
	router chi.Router
}

func NewServerWithRouter(server HTTPServer, router chi.Router) *ServerWithRouter {
	return &ServerWithRouter{
		server: server,
		router: router,
	}
}

func (s *ServerWithRouter) AddRouter(router chi.Router) {
	s.router.Mount("/", router)
	s.server.SetHandler(s.router)
}

func (s *ServerWithRouter) Run(ctx context.Context) error {
	s.start(ctx)
	<-ctx.Done()
	return s.stop(ctx)
}

func (s *ServerWithRouter) start(ctx context.Context) error {
	log.Info("Starting server...")
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("Server started with error", "error", err)
		return err
	}
	return nil
}

func (s *ServerWithRouter) stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := s.server.Shutdown(ctx)
	if err != nil {
		log.Error("Server stopped with error", "error", err)
		return err
	}
	log.Info("Server gracefully stopped")
	return nil
}
