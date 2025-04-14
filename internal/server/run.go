package server

import (
	"context"
	"net/http"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/log"
)

func Run(ctx context.Context, srv *http.Server) error {
	log.Info("Starting server...")

	go func() {
		srv.ListenAndServe()
	}()

	log.Info("Server is waiting for shutdown signal...")
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Info("Server is shutting down...")

	return srv.Shutdown(shutdownCtx)
}
