package middlewares

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"

	"go.uber.org/zap"
)

func TxMiddleware(
	db *sqlx.DB,
	txSetter func(ctx context.Context, tx *sqlx.Tx) context.Context,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if logger.Logger != nil {
				logger.Logger.Info("Starting a new transaction")
			}

			var tx *sqlx.Tx
			var err error
			if db != nil {
				tx, err = db.Beginx()
				if err != nil {
					if logger.Logger != nil {
						logger.Logger.Error("Failed to start transaction", zap.Error(err))
					}
					handleError(w)
					return
				}

			}

			var ctx context.Context
			if tx != nil {
				ctx = txSetter(r.Context(), tx)
			} else {
				ctx = r.Context()
			}

			rw := newBufferedResponseWriter(w)

			next.ServeHTTP(rw, r.WithContext(ctx))

			if logger.Logger != nil {
				logger.Logger.Info("Request completed", zap.Int("status", rw.status))
			}

			if tx != nil {
				if rw.status >= 400 {
					if logger.Logger != nil {
						logger.Logger.Warn("Rolling back transaction due to client or server error", zap.Int("status", rw.status))
					}
					if rollbackErr := tx.Rollback(); rollbackErr != nil {
						if logger.Logger != nil {
							logger.Logger.Error("Failed to rollback transaction", zap.Error(rollbackErr))
						}
						handleError(w)
						return
					}
					if logger.Logger != nil {
						logger.Logger.Info("Transaction rolled back successfully")
					}
					rw.FlushToUnderlying()
				} else {
					if logger.Logger != nil {
						logger.Logger.Info("Committing transaction")
					}
					if commitErr := tx.Commit(); commitErr != nil {
						if logger.Logger != nil {
							logger.Logger.Error("Failed to commit transaction", zap.Error(commitErr))
						}
						handleError(w)
						return
					}
					if logger.Logger != nil {
						logger.Logger.Info("Transaction committed successfully")
					}
					rw.FlushToUnderlying()
				}
			}
		})
	}
}

func handleError(w http.ResponseWriter) {
	http.Error(w, "Internal  server error", http.StatusInternalServerError)
}

type contextTxKey string

const txKey contextTxKey = "tx"

func SetTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func GetTx(ctx context.Context) *sqlx.Tx {
	tx, _ := ctx.Value(txKey).(*sqlx.Tx)
	return tx
}

type bufferedResponseWriter struct {
	http.ResponseWriter
	status  int
	headers http.Header
	body    []byte
	wrote   bool
}

func newBufferedResponseWriter(w http.ResponseWriter) *bufferedResponseWriter {
	return &bufferedResponseWriter{
		ResponseWriter: w,
		headers:        make(http.Header),
	}
}

func (rw *bufferedResponseWriter) Header() http.Header {
	return rw.headers
}

func (rw *bufferedResponseWriter) WriteHeader(statusCode int) {
	if rw.wrote {
		return
	}
	rw.status = statusCode
	rw.wrote = true
}

func (rw *bufferedResponseWriter) Write(b []byte) (int, error) {
	if !rw.wrote {
		rw.WriteHeader(http.StatusOK)
	}
	rw.body = append(rw.body, b...)
	return len(b), nil
}

func (rw *bufferedResponseWriter) FlushToUnderlying() {
	for k, vv := range rw.headers {
		for _, v := range vv {
			rw.ResponseWriter.Header().Add(k, v)
		}
	}
	rw.ResponseWriter.WriteHeader(rw.status)
	_, _ = rw.ResponseWriter.Write(rw.body)
}
