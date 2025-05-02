package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func TxMiddleware(
	db *sqlx.DB,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tx, err := db.BeginTxx(r.Context(), nil)
			if err != nil {
				return
			}

			defer func() {
				if rec := recover(); rec != nil {
					tx.Rollback()
					panic(rec)
				} else {
					if err := tx.Commit(); err != nil {
						tx.Rollback()
					}
				}
			}()

			ctx := setTx(r.Context(), tx)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

type txContextKey struct{}

func GetTxFromContext(ctx context.Context) (*sqlx.Tx, error) {
	tx, ok := ctx.Value(txContextKey{}).(*sqlx.Tx)
	if !ok {
		return nil, errors.New("transaction not found in context")
	}
	return tx, nil
}

func setTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txContextKey{}, tx)
}
