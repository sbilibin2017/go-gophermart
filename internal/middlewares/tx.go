package middlewares

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/db"
)

// TxMiddleware manages database transactions within the request lifecycle
func TxMiddleware(d *sqlx.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Using WithTx to handle transaction lifecycle
			err := db.WithTx(r.Context(), d, func(tx *sqlx.Tx) error {
				// Store the transaction in the context for downstream handlers
				ctxWithTx := context.WithValue(r.Context(), db.TxKey, tx)
				next.ServeHTTP(rw, r.WithContext(ctxWithTx))
				return nil
			})

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// If everything is fine, just write the response
			rw.ResponseWriter.WriteHeader(rw.statusCode)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
