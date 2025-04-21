package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/apps"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/server"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
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
	runAddress  string
	databaseURI string
)

func flags() {
	flag.StringVar(&runAddress, "a", "", "run address")
	flag.StringVar(&databaseURI, "d", "", "database uri")
	flag.Parse()

	if envA := os.Getenv("RUN_ADDRESS"); envA != "" {
		runAddress = envA
	}
	if envD := os.Getenv("DATABASE_URI"); envD != "" {
		databaseURI = envD
	}
}

func run() error {
	logger.Logger.Info("Connecting to database", zap.String("database_uri", databaseURI))

	db, err := sqlx.Connect("pgx", databaseURI)
	if err != nil {
		logger.Logger.Fatal("Failed to open database", zap.Error(err))
		return err
	}
	defer db.Close()

	logger.Logger.Info("Successfully connected to database")

	logger.Logger.Info("Starting server", zap.String("address", runAddress))

	srv := &http.Server{Addr: runAddress}

	apps.InitializeAccrualApp(srv, db)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err = server.Run(ctx, srv)
	if err != nil {
		return err
	}

	return nil
}
