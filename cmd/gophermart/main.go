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
	runAddress           string
	databaseURI          string
	accrualSystemAddress string
)

func flags() {
	flag.StringVar(&runAddress, "a", "", "run address")
	flag.StringVar(&databaseURI, "d", "", "database uri")
	flag.StringVar(&accrualSystemAddress, "r", "", "accrual system address")

	flag.Parse()

	if envAddr := os.Getenv("RUN_ADDRESS"); envAddr != "" {
		runAddress = envAddr
	}
	if envDBURI := os.Getenv("DATABASE_URI"); envDBURI != "" {
		databaseURI = envDBURI
	}
	if envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		accrualSystemAddress = envAccrualSystemAddress
	}
}

func run() {
	logger.Logger.Infow("Starting gophermart",
		"address", runAddress,
		"databaseURI", databaseURI,
		"accrualSystemAddress", accrualSystemAddress,
	)

	db, err := sql.Open("pgx", databaseURI)
	if err != nil {
		logger.Logger.Errorw("Failed to connect to database", "error", err)
		return
	}
	defer db.Close()

	router := chi.NewRouter()

	server := &http.Server{
		Addr:    runAddress,
		Handler: router,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		logger.Logger.Infow("Starting HTTP server", "address", runAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Errorw("HTTP server error", "error", err)
			stop() // Завершим контекст в случае ошибки сервера
		}
	}()

	<-ctx.Done()
	logger.Logger.Infow("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Logger.Errorw("Error during server shutdown", "error", err)
		return
	}

	logger.Logger.Infow("Server shutdown gracefully")
}
