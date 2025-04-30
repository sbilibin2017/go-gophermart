package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/routers"

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

	router := chi.NewRouter()

	routers.RegisterAccrualRouter(
		router,
		db,
		"/api",
		nil,
		nil,
		nil,
		middlewares.LoggingMiddleware,
		middlewares.GzipMiddleware,
		middlewares.TxMiddleware(db, contextutils.SetTx),
	)

	return nil
}
