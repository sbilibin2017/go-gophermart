package server

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var serverLog *zap.Logger

func init() {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	serverLog, _ = zapConfig.Build()
}

type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

func RunWithGracefulShutdown(ctx context.Context, srv Server) error {
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverLog.Error("Server failed", zap.Error(err))
		}
	}()

	<-ctx.Done()
	serverLog.Info("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		serverLog.Error("Server shutdown failed", zap.Error(err))
		return err
	}

	serverLog.Info("Server gracefully stopped")
	return nil
}
