package storage

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger" // Импорт пакета логирования
)

func WithTx(
	ctx context.Context,
	db *sqlx.DB,
	op func(ctx context.Context, tx *sqlx.Tx) error,
) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Logger.Errorw("Failed to begin transaction", "error", err)
		return err
	}
	if err := op(ctx, tx); err != nil {
		logger.Logger.Errorw("Transaction operation failed", "error", err)
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		logger.Logger.Errorw("Failed to commit transaction", "error", err)
		return err
	}
	return nil
}
