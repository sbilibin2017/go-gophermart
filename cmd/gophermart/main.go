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
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
)

func main() {
	flags()
	run()
}

type Config struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
	JWTSecretKey         []byte
	JWTExp               time.Duration
}

var config = Config{
	JWTSecretKey: []byte("test"),
	JWTExp:       time.Duration(365 * 24 * time.Hour),
}

func flags() {
	flag.StringVar(&config.RunAddress, "a", "", "run address")
	flag.StringVar(&config.DatabaseURI, "d", "", "database uri")
	flag.StringVar(&config.AccrualSystemAddress, "r", "", "accrual system address")

	flag.Parse()

	if envA := os.Getenv("RUN_ADDRESS"); envA != "" {
		config.RunAddress = envA
	}
	if envD := os.Getenv("DATABASE_URI"); envD != "" {
		config.DatabaseURI = envD
	}
	if envR := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envR != "" {
		config.AccrualSystemAddress = envR
	}
}

func run() {
	logger.Logger.Infof("Starting gophermart on address: %s", config.RunAddress)
	logger.Logger.Infof("Database URI: %s", config.DatabaseURI)
	logger.Logger.Infof("Accrual System Address: %s", config.AccrualSystemAddress)

	db, err := sql.Open("pgx", config.DatabaseURI)
	if err != nil {
		logger.Logger.Errorf("Failed to connect to database: %v", err)
		return
	}
	err = db.Ping()
	if err != nil {
		logger.Logger.Errorf("Failed to ping database: %v", err)
		return
	}
	defer db.Close()

	rtr := chi.NewRouter()

	rtr.Use(
		middlewares.LoggingMiddleware,
		middlewares.GzipMiddleware,
		middlewares.TxMiddleware(db),
	)

	rtr.Route("/api/user", func(api chi.Router) {
		api.Post("/register", nil)
		api.Post("/login", nil)
		api.With(middlewares.AuthMiddleware(config.JWTSecretKey)).Post("/orders", nil)
		api.With(middlewares.AuthMiddleware(config.JWTSecretKey)).Get("/orders", nil)
		api.With(middlewares.AuthMiddleware(config.JWTSecretKey)).Get("/balance", nil)
		api.With(middlewares.AuthMiddleware(config.JWTSecretKey)).Post("/balance/withdraw", nil)
		api.With(middlewares.AuthMiddleware(config.JWTSecretKey)).Get("/withdrawals", nil)
	})

	srv := &http.Server{Addr: config.RunAddress, Handler: rtr}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Errorf("Server error: %v", err)
		}
	}()

	logger.Logger.Info("Server is running... Press Ctrl+C to stop.")

	<-ctx.Done()
	logger.Logger.Info("Shutdown signal received, stopping server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Logger.Errorf("Error during shutdown: %v", err)
	} else {
		logger.Logger.Info("Server gracefully stopped.")
	}
}
