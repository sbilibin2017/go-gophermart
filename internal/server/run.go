package server

import (
	"context"
	"net/http"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/log"
)

type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

func Run(ctx context.Context, srv Server) error {
	log.Info("Starting server...")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("Server started with error", "error", err)
		return err
	}
	<-ctx.Done()
	log.Info("Stopping server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Error("Server stopped with error", "error", err)
		return err
	}
	log.Info("Server gracefully stopped")
	return nil
}
