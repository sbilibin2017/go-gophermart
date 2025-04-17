package db

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var opener = sqlx.Open

func NewDB(dsn string) (*sqlx.DB, error) {
	db, err := opener("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
