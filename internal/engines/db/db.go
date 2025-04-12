package db

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DatabaseURIGetter interface {
	GetDatabaseURI() string
}

func NewDB(d DatabaseURIGetter) *sql.DB {
	db, _ := sql.Open("pgx", d.GetDatabaseURI())
	return db
}
