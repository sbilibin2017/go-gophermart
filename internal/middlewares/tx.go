package middlewares

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"go.uber.org/zap"
)

func TxMiddleware(db *sqlx.DB, txContextSetter func(ctx context.Context, tx *sqlx.Tx) context.Context) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tx, err := db.BeginTxx(r.Context(), nil)
			if err != nil {
				logger.Logger.Error("Failed to begin transaction", zap.Error(err))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			r = r.WithContext(txContextSetter(r.Context(), tx))
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(rw, r)

			if tx == nil {
				return
			}
			if rw.statusCode >= http.StatusBadRequest {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					logger.Logger.Error("Failed to rollback transaction", zap.Error(rollbackErr))
				} else {
					logger.Logger.Info("Transaction rolled back due to error response")
				}
				return
			}

			if commitErr := tx.Commit(); commitErr != nil {
				logger.Logger.Error("Failed to commit transaction", zap.Error(commitErr))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
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
