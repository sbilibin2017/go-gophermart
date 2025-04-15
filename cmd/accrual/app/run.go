package app

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger" // Импортируй свой логгер
	"go.uber.org/zap/zapcore"
)

const (
	DatabaseDriver = "pgx"
)

func Run() {
	logger.Init(zapcore.InfoLevel)

	logger.Logger.Infow("Starting application")

	db, err := sqlx.Connect(DatabaseDriver, flagDatabaseURI)
	if err != nil {
		logger.Logger.Errorw("Failed to connect to database", "error", err)
		return
	}
	logger.Logger.Infow("Database connection opened")

	defer func() {
		if err := db.Close(); err != nil {
			logger.Logger.Warnw("Error closing database", "error", err)
		} else {
			logger.Logger.Infow("Database connection closed")
		}
	}()

	srv := &http.Server{Addr: flagRunAddr}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		logger.Logger.Infow("Server is starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Errorw("Server error", "error", err)
		}
	}()

	<-ctx.Done()
	logger.Logger.Infow("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Logger.Errorw("Server shutdown failed", "error", err)
		return
	}

	logger.Logger.Infow("Server shutdown completed gracefully")
}
