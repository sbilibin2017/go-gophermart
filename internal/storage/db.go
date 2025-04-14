package storage

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type DatabaseURIGetter interface {
	GetDatabaseURI() string
}

func NewDB(g DatabaseURIGetter) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", g.GetDatabaseURI())
	if err != nil {
		return nil, err
	}
	return db, err
}
