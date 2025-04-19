package middlewares

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type txKeyType string

const TxKey txKeyType = "tx"

func TxFromContext(ctx context.Context) *sqlx.Tx {
	tx, _ := ctx.Value(TxKey).(*sqlx.Tx)
	return tx
}

func TxMiddleware(db *sqlx.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			tx, err := db.BeginTxx(r.Context(), nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			defer func() {
				if rw.statusCode >= 400 {
					tx.Rollback()
				} else {
					tx.Commit()
				}
			}()

			ctxWithTx := context.WithValue(r.Context(), TxKey, tx)
			next.ServeHTTP(rw, r.WithContext(ctxWithTx))
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
