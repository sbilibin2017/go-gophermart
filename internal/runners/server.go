package runners

import (
	"context"
	"net/http"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

func RunServer(ctx context.Context, srv *http.Server) error {
	logger.Logger.Info("Starting server...")

	go func() {
		srv.ListenAndServe()
	}()

	logger.Logger.Info("Server is waiting for shutdown signal...")
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Logger.Info("Server is shutting down...")

	return srv.Shutdown(shutdownCtx)
}
