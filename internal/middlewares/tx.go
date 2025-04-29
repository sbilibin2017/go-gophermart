package middlewares

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/logger" // Импорт пакета логирования
)

func TxMiddleware(db *sqlx.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := withTx(r.Context(), db, func(ctx context.Context, tx *sqlx.Tx) error {
				r = r.WithContext(contextutils.SetTx(ctx, tx))
				next.ServeHTTP(w, r)
				return nil
			})
			if err != nil {
				logger.Logger.Errorw("Transaction failed", "error", err)
			}
		})
	}
}

func withTx(
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
