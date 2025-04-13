package storage

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Transaction struct {
	db *sqlx.DB
}

func NewTransaction(db *sqlx.DB) *Transaction {
	return &Transaction{db: db}
}

func (t *Transaction) Do(ctx context.Context, operation func(tx *sql.Tx) error) error {
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
