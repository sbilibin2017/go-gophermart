package middlewares

import (
	"context"
	"database/sql"
	"log"
	"net/http"
)

type txKeyType string

const txKey txKeyType = "tx"

func TxFromContext(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey).(*sql.Tx)
	return tx, ok
}

func TxMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tx, err := db.Begin()
			if err != nil {
				http.Error(w, "failed to begin transaction", http.StatusInternalServerError)
				return
			}

			ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			ctx := context.WithValue(r.Context(), txKey, tx)
			r = r.WithContext(ctx)

			next.ServeHTTP(ww, r)

			if ww.statusCode >= http.StatusBadRequest {
				log.Printf("Rolling back transaction due to status code: %d", ww.statusCode)
				tx.Rollback()
				return
			}

			if err := tx.Commit(); err != nil {
				log.Printf("Failed to commit transaction: %v", err)
				http.Error(w, "transaction commit error", http.StatusInternalServerError)
				return
			}
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
