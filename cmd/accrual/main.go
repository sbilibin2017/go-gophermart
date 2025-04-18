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
	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
)

func main() {
	flags()
	run()
}

type Config struct {
	RunAddress  string
	DatabaseURI string
}

var config Config

func flags() {
	flag.StringVar(&config.RunAddress, "a", "", "run address")
	flag.StringVar(&config.DatabaseURI, "d", "", "database uri")

	flag.Parse()

	if envA := os.Getenv("RUN_ADDRESS"); envA != "" {
		config.RunAddress = envA
	}
	if envD := os.Getenv("DATABASE_URI"); envD != "" {
		config.DatabaseURI = envD
	}
}

func run() {
	logger.Logger.Infof("Starting gophermart on address: %s", config.RunAddress)
	logger.Logger.Infof("Database URI: %s", config.DatabaseURI)

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

	val := validator.New()

	mws := []func(http.Handler) http.Handler{
		middlewares.LoggingMiddleware,
		middlewares.GzipMiddleware,
		middlewares.TxMiddleware(db),
	}

	rtr.Route("/api", func(api chi.Router) {
		api.Use(mws...)

		api.Get("/orders/{number}", nil)
		api.Post("/orders", nil)
		api.Post("/goods", handlers.RegisterGoodRewardHandler(val, nil))
	})

	srv := &http.Server{Addr: config.RunAddress, Handler: rtr}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Errorf("Server error: %v", err)
		}
	}()

	logger.Logger.Info("Server is running... Press Ctrl+C to stop.")

	<-ctx.Done()
	logger.Logger.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Logger.Errorf("Server shutdown failed: %v", err)
	} else {
		logger.Logger.Info("Server gracefully stopped.")
	}
}
