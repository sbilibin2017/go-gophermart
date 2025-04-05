package server

import (
	"context"
	"net/http"

	c "github.com/sbilibin2017/go-gophermart/pkg/context"
	"github.com/sbilibin2017/go-gophermart/pkg/log"
)

type HTTPServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

type ServerWithContext struct {
	srv HTTPServer
}

func NewServerWithContext(srv HTTPServer) *ServerWithContext {
	return &ServerWithContext{
		srv: srv,
	}
}

func (s *ServerWithContext) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		log.Info("Starting server on")
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()
	select {
	case <-ctx.Done():
		log.Info("Context canceled")
	case err := <-errCh:
		log.Error("Server error", "error", err)
		return err
	}
	shutdownCtx, cancel := c.NewShutdownContext()
	defer cancel()
	log.Info("Shutting down server gracefully...")
	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		log.Error("Shutdown error", "error", err)
		return err
	}
	log.Info("Server shutdown completed")
	return nil
}
