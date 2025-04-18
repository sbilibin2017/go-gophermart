package storage

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var connOpener func(driverName, dataSourceName string) (*sql.DB, error) = sql.Open

func NewDB(dsn string) (*sql.DB, error) {
	db, err := connOpener("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
