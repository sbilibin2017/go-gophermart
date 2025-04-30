package middlewares

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func TxMiddleware(
	db *sqlx.DB,
	txSetter func(ctx context.Context, tx *sqlx.Tx) context.Context,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tx, err := db.BeginTxx(r.Context(), nil)
			if err != nil {
				http.Error(w, "failed to start transaction: "+err.Error(), http.StatusInternalServerError)
				return
			}

			defer func() {
				if rec := recover(); rec != nil {
					tx.Rollback()
					panic(rec)
				} else {
					if err := tx.Commit(); err != nil {
						tx.Rollback()
						http.Error(w, "Internal server error", http.StatusInternalServerError)
					}
				}
			}()

			ctx := txSetter(r.Context(), tx)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
