package unitofwork

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type UnitOfWork struct {
	db *sqlx.DB
}

func NewUnitOfWork(db *sqlx.DB) *UnitOfWork {
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
