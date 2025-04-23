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

func RunWithGracefulShutdown(ctx context.Context, srv Server) error {
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Error("Server failed", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Logger.Info("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Logger.Error("Server shutdown failed", zap.Error(err))
		return err
	}

	logger.Logger.Info("Server gracefully stopped")
	return nil
}
