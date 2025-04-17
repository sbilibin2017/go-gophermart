package main

import (
	"context"
	"database/sql"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"go.uber.org/zap/zapcore"
)

func main() {
	flags()
	run()
}

func init() {
	logger.Init(zapcore.InfoLevel)
}

var (
	a string
	d string
	r string
)

func flags() {
	flag.StringVar(&a, "a", "", "run address")
	flag.StringVar(&d, "d", "", "database uri")
	flag.StringVar(&r, "r", "", "accrual system address")

	flag.Parse()

	if envA := os.Getenv("RUN_ADDRESS"); envA != "" {
		a = envA
	}
	if envD := os.Getenv("DATABASE_URI"); envD != "" {
		d = envD
	}
	if envR := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envR != "" {
		r = envR
	}
}

func run() {
	logger.Logger.Infow(
		"Starting gophermart",
		"address", a,
		"databaseURI", d,
		"accrualSystemAddress", r,
	)

	db, err := sql.Open("pgx", d)
	if err != nil {
		logger.Logger.Errorw("Failed to connect to database", "error", err)
		return
	}
	defer db.Close()

	router := chi.NewRouter()

	server := &http.Server{
		Addr:    a,
		Handler: router,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		logger.Logger.Infow("Starting HTTP server", "address", a)
		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			logger.Logger.Errorw("HTTP server error", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	logger.Logger.Infow("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Logger.Errorw("Error during server shutdown", "error", err)
		return
	}

	logger.Logger.Infow("Server shutdown gracefully")
}
