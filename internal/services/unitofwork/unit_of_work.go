package unitofwork

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type UnitOfWork struct {
	db *sql.DB
}

func NewUnitOfWork(db *sql.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

func (uow *UnitOfWork) Do(ctx context.Context, operation func(tx *sql.Tx) error) error {
	if uow.db == nil {
		return nil
	}
	tx, _ := uow.db.BeginTx(ctx, nil)

	err := operation(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
