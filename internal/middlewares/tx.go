package middlewares

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func TxMiddleware(
	conn *sqlx.DB,
	ctxProvider func(ctx context.Context, tx *sqlx.Tx) context.Context,
	txProvider func(ctx context.Context, db *sqlx.DB, fn func(tx *sqlx.Tx) error) error,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			err := txProvider(ctx, conn, func(tx *sqlx.Tx) error {
				ctxWithTx := ctxProvider(ctx, tx)
				r = r.WithContext(ctxWithTx)
				next.ServeHTTP(w, r)
				return nil
			})

			if err != nil {
				http.Error(w, "Internal error", http.StatusInternalServerError)
			}
		})
	}
}
