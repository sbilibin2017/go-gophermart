package db

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

var opener = sqlx.Open

func NewDB(dsn string) (*sqlx.DB, error) {
	logger.Logger.Info("Попытка подключения к базе данных с DSN:", dsn)
	db, err := opener("pgx", dsn)
	if err != nil {
		logger.Logger.Error("Не удалось подключиться к базе данных:", err)
		return nil, err
	}
	logger.Logger.Info("Успешное подключение к базе данных.")
	return db, nil
}
