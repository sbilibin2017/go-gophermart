package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/apps"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/server"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func run(config *configs.AccrualConfig) error {
	logger.InitLoggerWithInfoLevel()
	defer logger.Logger.Sync()

	db, err := sqlx.Connect("pgx", config.DatabaseURI)
	if err != nil {
		return err
	}
	defer db.Close()

	srv := &http.Server{Addr: config.RunAddress, Handler: chi.NewRouter()}

	apps.ConfigureAccrualApp(db, srv)

	return server.Run(
		context.Background(),
		srv,
	)
}
