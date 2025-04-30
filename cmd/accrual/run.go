package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/logger"

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

}
