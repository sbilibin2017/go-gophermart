package middlewares

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
)

// TxMiddleware - middleware для работы с транзакцией
func TxMiddleware(db *sqlx.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Начинаем новую транзакцию
			tx, err := db.Beginx()
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Передаем транзакцию через контекст запроса, используя custom contextKey type
			ctx := context.WithValue(r.Context(), txKey, tx)
			r = r.WithContext(ctx)

			// Используем defer для отката транзакции, если произойдет ошибка
			defer func() {
				if err != nil {
					// Если произошла ошибка, откатываем транзакцию
					if rollbackErr := tx.Rollback(); rollbackErr != nil {
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					}
				}
			}()

			// Завершаем выполнение запроса
			next.ServeHTTP(w, r)

			// Если код завершился без ошибок, коммитим транзакцию
			if err := tx.Commit(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		})
	}
}

func GetTxFromContext(ctx context.Context) *sqlx.Tx {
	tx, ok := ctx.Value(txKey).(*sqlx.Tx)
	if !ok {
		return nil
	}
	return tx
}

// Define a custom type for the context key to avoid using the built-in string type.
type contextKey string

// Define a constant for the context key.
const txKey contextKey = "tx"
