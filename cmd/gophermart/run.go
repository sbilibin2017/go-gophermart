package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"

	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func run() error {
	db, err := sqlx.Connect("pgx", options.databaseURI)
	if err != nil {
		logger.Error("Error connecting to the database", "error", err)
		return err
	}
	defer db.Close()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	srv := &http.Server{Addr: options.runAddress}

	go func() {
		logger.Info("Starting server on " + options.runAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Error while starting the server", "error", err)
		}
	}()

	<-ctx.Done()
	logger.Info("Received shutdown signal")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Error during server shutdown", "error", err)
	} else {
		logger.Info("Server shutdown gracefully")
	}

	return nil
}
