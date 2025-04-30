package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

func getExecutor(
	ctx context.Context,
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
) sqlx.ExtContext {
	tx, err := txProvider(ctx)
	if err != nil {
		return db
	}
	return tx
}
