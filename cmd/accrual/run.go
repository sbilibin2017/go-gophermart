package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/apps"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/server"
	"go.uber.org/zap/zapcore"
)

func run(config *configs.AccrualConfig) error {
	err := logger.Init(zapcore.InfoLevel)
	if err != nil {
		return err
	}

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
		os.Exit(1)
	}
	os.Exit(0)
}
