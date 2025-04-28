package middlewares

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

func TxMiddleware(
	db *sqlx.DB,
	txSetter func(ctx context.Context, tx *sqlx.Tx) context.Context,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := withTx(db, func(tx *sqlx.Tx) error {
				ctx := txSetter(r.Context(), tx)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return nil
			})
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		})
	}
}

func withTx(db *sqlx.DB, op func(tx *sqlx.Tx) error) error {
	tx, err := db.Beginx()
	if err != nil {
		logger.Logger.Error("Error starting transaction: ", err)
		return err
	}
	logger.Logger.Info("Executing transaction operation")
	err = op(tx)
	if err != nil {
		logger.Logger.Error("Error during transaction operation: ", err)
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			logger.Logger.Error("Error rolling back transaction: ", rollbackErr)
			return rollbackErr
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		logger.Logger.Error("Error committing transaction: ", err)
		return err
	}
	logger.Logger.Info("Transaction successfully committed")
	return nil
}
