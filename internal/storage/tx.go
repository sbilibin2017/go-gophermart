package storage

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func WithTx(db *sql.DB, op func(tx *sql.Tx) error) error {
	if db == nil {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = op(tx)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return tx.Commit()
}
