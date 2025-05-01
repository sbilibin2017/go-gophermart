package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

func Run(ctx context.Context, srv *http.Server) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	errChan := make(chan error, 1)

	go func() {
		logger.Logger.Info("Starting server...")
		errChan <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		logger.Logger.Info("Shutting down server...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return err
		}
		logger.Logger.Info("Server gracefully stopped")
		return nil

	case err := <-errChan:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	}
}
