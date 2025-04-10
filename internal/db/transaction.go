package db

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Transaction struct {
	db *sql.DB
}

func NewTransaction(db *sql.DB) *Transaction {
	return &Transaction{db: db}
}

func (t *Transaction) Do(ctx context.Context, operation func(tx *sql.Tx) error) error {
	var err error
	var tx *sql.Tx
	tx, err = t.db.BeginTx(ctx, nil)
	if err != nil {
		return ErrFailedToBeginTransaction
	}
	err = operation(tx)
	if err != nil {
		tx.Rollback()
		return ErrFailedToRollbackTransaction
	}
	if err := tx.Commit(); err != nil {
		return ErrFailedToCommitTransaction
	}
	return nil
}

var (
	ErrFailedToBeginTransaction    = errors.New("failed to begin transaction")
	ErrFailedToCommitTransaction   = errors.New("failed to commit transaction")
	ErrFailedToRollbackTransaction = errors.New("failed to rollback transaction")
)
