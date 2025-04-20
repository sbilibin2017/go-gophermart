package db

import (
	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var connFactory = sqlx.Connect

func NewDB(dsn string) (*sqlx.DB, error) {
	db, err := connFactory("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return db, err
}
