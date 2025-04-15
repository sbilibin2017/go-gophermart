package storage

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/log" // Замените на актуальный путь к вашему пакету log
)

var connect = sqlx.Connect

func NewDB(dsn string) (*sqlx.DB, error) {
	log.Info("Попытка подключения к базе данных", "dsn", dsn)
	db, err := connect("pgx", dsn)
	if err != nil {
		log.Error("Ошибка при подключении к базе данных", "ошибка", err)
		return nil, err
	}
	log.Info("Успешное подключение к базе данных")
	return db, nil
}
