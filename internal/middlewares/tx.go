package middlewares

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/handlers/utils"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"go.uber.org/zap"
)

type txKeyType string

const txKey txKeyType = "tx"

func TxFromContext(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		return tx
	}
	return nil
}

func TxMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tx, err := beginTx(db, w, r)
			if err != nil {
				return
			}
			ctxWithTx := context.WithValue(r.Context(), txKey, tx)
			r = r.WithContext(ctxWithTx)
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(rw, r)
			handleTransactionEnd(tx, rw.statusCode, w)
		})
	}
}

func beginTx(db *sql.DB, w http.ResponseWriter, r *http.Request) (*sql.Tx, error) {
	tx, err := db.BeginTx(r.Context(), nil)
	if err != nil {
		logger.Logger.Error("Failed to begin transaction", zap.Error(err))
		utils.ErrorInternalServerResponse(w, err)
		return nil, err
	}
	return tx, nil
}

func handleTransactionEnd(tx *sql.Tx, statusCode int, w http.ResponseWriter) {
	if tx == nil {
		return
	}
	if statusCode >= http.StatusBadRequest {
		handleRollback(tx)
		return
	}
	handleCommit(tx, w)
}

func handleRollback(tx *sql.Tx) {
	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		logger.Logger.Error("Failed to rollback transaction", zap.Error(rollbackErr))
	} else {
		logger.Logger.Info("Transaction rolled back due to error response")
	}
}

func handleCommit(tx *sql.Tx, w http.ResponseWriter) {
	commitErr := tx.Commit()
	if commitErr != nil {
		logger.Logger.Error("Failed to commit transaction", zap.Error(commitErr))
		utils.ErrorInternalServerResponse(w, commitErr)
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
