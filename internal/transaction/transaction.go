package transaction

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	ErrFailedToBeginTransaction    = errors.New("failed to begin transaction")
	ErrFailedToCommitTransaction   = errors.New("failed to commit transaction")
	ErrFailedToRollbackTransaction = errors.New("failed to rollback transaction")
)

type DB interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
}

type Tx interface {
	Commit() error
	Rollback() error
}

type Transaction struct {
	db DB
}

func NewTransaction(db DB) *Transaction {
	return &Transaction{db: db}
}

func (t *Transaction) Do(ctx context.Context, operation func(tx Tx) error) error {
	tx, err := t.db.BeginTx(ctx, nil)
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
