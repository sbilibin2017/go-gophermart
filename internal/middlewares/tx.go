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
			// Начинаем транзакцию с использованием контекста запроса
			tx, err := db.BeginTxx(r.Context(), nil)
			if err != nil {
				http.Error(w, "failed to start transaction: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// Обработка отката транзакции в случае паники или ошибки
			defer func() {
				if rec := recover(); rec != nil {
					tx.Rollback()
					panic(rec) // передаем панику дальше после отката
				} else {
					if err := tx.Commit(); err != nil {
						// Если не удалось зафиксировать транзакцию, откатываем
						tx.Rollback()
						http.Error(w, "failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
					}
				}
			}()

			// Устанавливаем транзакцию в контекст запроса
			ctx := txSetter(r.Context(), tx)
			r = r.WithContext(ctx)

			// Переход к следующему обработчику
			next.ServeHTTP(w, r)
		})
	}
}
