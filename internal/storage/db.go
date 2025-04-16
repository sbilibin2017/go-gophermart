package storage

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var connOpener = sqlx.Open

func NewDB(dsn string) (*sqlx.DB, error) {
	db, err := connOpener("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
