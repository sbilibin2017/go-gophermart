package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/apps"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/server"
	"go.uber.org/zap"
)

func main() {
	config := flags()
	err := run(config)
	exit(err)
}

var (
	runAddress  string
	databaseURI string
)

func flags() *configs.AccrualConfig {

	flag.StringVar(&runAddress, "a", "", "run address")
	flag.StringVar(&databaseURI, "d", "", "database uri")
	flag.Parse()

	if envA := os.Getenv("RUN_ADDRESS"); envA != "" {
		runAddress = envA
	}
	if envD := os.Getenv("DATABASE_URI"); envD != "" {
		databaseURI = envD
	}

	return &configs.AccrualConfig{
		RunAddress:  runAddress,
		DatabaseURI: databaseURI,
	}
}

func run(config *configs.AccrualConfig) error {
	logger.Init()

	db, err := sqlx.Connect("pgx", config.DatabaseURI)
	if err != nil {
		return err
	}

	defer func() {
		db.Close()
		logger.Logger.Sync()
	}()

	srv := &http.Server{
		Addr: config.RunAddress,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	apps.ConfigureAccrualServer(db, srv)

	err = server.RunWithGracefulShutdown(ctx, srv)
	if err != nil {
		return err
	}

	return nil
}

func exit(err error) {
	if err != nil {
		logger.Logger.Error("Error starting accrual app",
			zap.Error(err),
		)
		os.Exit(1)
	}
	os.Exit(0)
}
