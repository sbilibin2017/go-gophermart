package server

import (
	"context"
	"net/http"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

func Run(ctx context.Context, srv *http.Server) error {
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Errorf("Server failed: %v", err)
		}
	}()

	logger.Logger.Infof("Server started successfully")

	<-ctx.Done()
	sigReceived := ctx.Err()
	logger.Logger.Infof("Received signal: %v. Initiating shutdown...", sigReceived)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Logger.Errorf("Server shutdown failed: %v", err)
		return err
	}

	logger.Logger.Info("Server stopped gracefully.")
	return nil
}
