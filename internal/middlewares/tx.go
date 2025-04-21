package middlewares

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	hh "github.com/sbilibin2017/go-gophermart/internal/handlers/helpers"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	rh "github.com/sbilibin2017/go-gophermart/internal/repositories/helpers"
	"go.uber.org/zap"
)

func TxMiddleware(db *sqlx.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tx, err := beginTx(db, w, r)
			if err != nil {
				logger.Logger.Error("Failed to begin transaction", zap.Error(err))
				return
			}
			ctxWithTx := rh.TxToContext(r.Context(), tx)
			r = r.WithContext(ctxWithTx)
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(rw, r)
			handleTransactionEnd(tx, rw.statusCode, w)
		})
	}
}

func beginTx(db *sqlx.DB, w http.ResponseWriter, r *http.Request) (*sqlx.Tx, error) {
	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		logger.Logger.Error("Failed to begin transaction", zap.Error(err))
		hh.ErrorInternalServerResponse(w, err)
		return nil, err
	}
	return tx, nil
}

func handleTransactionEnd(tx *sqlx.Tx, statusCode int, w http.ResponseWriter) {
	if tx == nil {
		return
	}
	if statusCode >= http.StatusBadRequest {
		handleRollback(tx)
		return
	}
	handleCommit(tx, w)
}

func handleRollback(tx *sqlx.Tx) {
	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		logger.Logger.Error("Failed to rollback transaction", zap.Error(rollbackErr))
	} else {
		logger.Logger.Info("Transaction rolled back due to error response")
	}
}

func handleCommit(tx *sqlx.Tx, w http.ResponseWriter) {
	commitErr := tx.Commit()
	if commitErr != nil {
		logger.Logger.Error("Failed to commit transaction", zap.Error(commitErr))
		hh.ErrorInternalServerResponse(w, commitErr)
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
