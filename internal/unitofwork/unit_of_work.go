package unitofwork

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type UnitOfWork struct {
	db *sql.DB
}

func NewUnitOfWork(db *sql.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

func (t *UnitOfWork) Do(ctx context.Context, operation func(tx *sql.Tx) error) error {
	if t.db == nil {
		return nil
	}
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return ErrFailedToBeginTransaction
	}
	if err := operation(tx); err != nil {
		_ = tx.Rollback()
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
