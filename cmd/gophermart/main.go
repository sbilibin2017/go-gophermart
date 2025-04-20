package main

import (
	"context"
	"database/sql"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/server"
)

func main() {
	flags()
	err := run()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

var (
	runAddress           string
	databaseURI          string
	accrualSystemAddress string
	jwtSecretKey         = []byte("test")
	// jwtExp               = time.Duration(365 * 24 * time.Hour)
)

func flags() {
	flag.StringVar(&runAddress, "a", "", "run address")
	flag.StringVar(&databaseURI, "d", "", "database uri")
	flag.StringVar(&accrualSystemAddress, "r", "", "accrual system address")

	flag.Parse()

	if envA := os.Getenv("RUN_ADDRESS"); envA != "" {
		runAddress = envA
	}
	if envD := os.Getenv("DATABASE_URI"); envD != "" {
		databaseURI = envD
	}
	if envR := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envR != "" {
		accrualSystemAddress = envR
	}
}

func run() error {
	logger.Logger.Info("Connecting to database", zap.String("database_uri", databaseURI))

	db, err := sql.Open("pgx", databaseURI)
	if err != nil {
		logger.Logger.Error("Failed to open database", zap.Error(err))
		return err
	}
	defer db.Close()

	logger.Logger.Info("Successfully connected to database")

	api := chi.NewRouter()

	api.Route("/api/user", func(r chi.Router) {
		r.Use(
			middlewares.LoggingMiddleware,
			middlewares.GzipMiddleware,
			middlewares.TxMiddleware(db),
		)

		r.Post("/register", nil)
		r.Post("/login", nil)

		r.Use(
			middlewares.AuthMiddleware(jwtSecretKey),
		)

		r.Post("/orders", nil)
		r.Get("/orders", nil)
		r.Get("/balance", nil)
		r.Post("/balance/withdraw", nil)
		r.Get("/withdrawals", nil)
	})

	logger.Logger.Info("Starting server", zap.String("address", runAddress))

	srv := &http.Server{Addr: runAddress}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err = server.Run(ctx, srv)
	if err != nil {
		return err
	}

	return nil
}
