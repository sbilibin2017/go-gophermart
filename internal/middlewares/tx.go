package middlewares

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/logger" // Импорт пакета логирования
	"github.com/sbilibin2017/go-gophermart/internal/storage"
)

func TxMiddleware(db *sqlx.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := storage.WithTx(r.Context(), db, func(ctx context.Context, tx *sqlx.Tx) error {
				r = r.WithContext(contextutils.SetTx(ctx, tx))
				next.ServeHTTP(w, r)
				return nil
			})
			if err != nil {
				logger.Logger.Errorw("Transaction failed", "error", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		})
	}
}
