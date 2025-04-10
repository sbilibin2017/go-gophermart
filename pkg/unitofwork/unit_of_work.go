package unitofwork

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type UnitOfWork struct {
	db *sql.DB
}

func NewDBUnitOfWork(db *sql.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

func (uow *UnitOfWork) Do(ctx context.Context, operation func(tx *sql.Tx) error) error {
	var err error
	var tx *sql.Tx
	tx, err = uow.db.BeginTx(ctx, nil)
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
