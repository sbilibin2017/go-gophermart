package db

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Tx struct {
	db *sqlx.DB
}

func NewTx(db *sqlx.DB) *Tx {
	return &Tx{db: db}
}

func (t *Tx) Do(ctx context.Context, operation func(tx *sql.Tx) error) error {
	if t.db == nil {
		return nil
	}
	tx, _ := t.db.BeginTx(ctx, nil)
	err := operation(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
