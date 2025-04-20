package server

import (
	"context"
	"net/http"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"go.uber.org/zap"
)

type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

func Run(ctx context.Context, srv Server) error {
	errCh := make(chan error, 1)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Error("Server failed", zap.Error(err))
			errCh <- err
			return
		}
		errCh <- nil
	}()

	<-ctx.Done()

	logger.Logger.Info("Received shutdown signal")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Logger.Error("Server shutdown failed", zap.Error(err))
		return err
	}

	logger.Logger.Info("Server stopped gracefully")

	return <-errCh
}
