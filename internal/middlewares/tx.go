// Package middlewares содержит HTTP middleware для работы с транзакциями в базе данных.
package middlewares

import (
	"database/sql"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/db"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"go.uber.org/zap"
)

// TxMiddleware возвращает middleware, который управляет транзакциями базы данных.
// Он обрабатывает транзакцию для каждого запроса, начиная с её создания, и в зависимости от результата запроса
// либо коммитит, либо откатывает транзакцию.
//
// Параметры:
//   - db: объект базы данных (*sql.DB), который используется для создания новой транзакции.
//
// Возвращает:
//   - func(http.Handler) http.Handler: Middleware, который оборачивает хендлер для управления транзакциями.
func TxMiddleware(dbConn *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Logger.Info("Starting a new transaction")

			// Начинаем новую транзакцию
			tx, err := dbConn.Begin()
			if err != nil {
				logger.Logger.Error("Failed to start transaction", zap.Error(err))
				http.Error(w, "could not start transaction", http.StatusInternalServerError)
				return
			}

			// Устанавливаем транзакцию в контексте запроса
			ctx := db.SetTx(r.Context(), tx)

			// Создаем буферизованный ResponseWriter для отслеживания статуса ответа
			rw := newBufferedResponseWriter(w)

			// Выполняем следующий хендлер
			next.ServeHTTP(rw, r.WithContext(ctx))

			logger.Logger.Info("Request completed", zap.Int("status", rw.status))

			// Если статус ошибки (>= 400), откатываем транзакцию
			if rw.status >= 400 {
				logger.Logger.Warn("Rolling back transaction due to client or server error", zap.Int("status", rw.status))
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					logger.Logger.Error("Failed to rollback transaction", zap.Error(rollbackErr))
					http.Error(w, "could not rollback transaction", http.StatusInternalServerError)
					return
				}
				logger.Logger.Info("Transaction rolled back successfully")
				rw.FlushToUnderlying()
			} else {
				// Если нет ошибки, коммитим транзакцию
				logger.Logger.Info("Committing transaction")
				if commitErr := tx.Commit(); commitErr != nil {
					logger.Logger.Error("Failed to commit transaction", zap.Error(commitErr))
					http.Error(w, "could not commit transaction", http.StatusInternalServerError)
					return
				}
				logger.Logger.Info("Transaction committed successfully")
				rw.FlushToUnderlying()
			}
		})
	}
}

// bufferedResponseWriter это структура, которая оборачивает http.ResponseWriter для отслеживания
// статуса и содержимого ответа перед отправкой клиенту.
type bufferedResponseWriter struct {
	http.ResponseWriter
	status  int
	headers http.Header
	body    []byte
	wrote   bool
}

// newBufferedResponseWriter создает новый экземпляр bufferedResponseWriter для указанного ResponseWriter.
func newBufferedResponseWriter(w http.ResponseWriter) *bufferedResponseWriter {
	return &bufferedResponseWriter{
		ResponseWriter: w,
		headers:        make(http.Header),
	}
}

// Header возвращает заголовки ответа.
func (rw *bufferedResponseWriter) Header() http.Header {
	return rw.headers
}

// WriteHeader устанавливает статус ответа.
func (rw *bufferedResponseWriter) WriteHeader(statusCode int) {
	if rw.wrote {
		return
	}
	rw.status = statusCode
	rw.wrote = true
}

// Write записывает данные в тело ответа.
func (rw *bufferedResponseWriter) Write(b []byte) (int, error) {
	if !rw.wrote {
		rw.WriteHeader(http.StatusOK)
	}
	rw.body = append(rw.body, b...)
	return len(b), nil
}

// FlushToUnderlying сбрасывает все буферизированные данные в оригинальный ResponseWriter.
func (rw *bufferedResponseWriter) FlushToUnderlying() {
	for k, vv := range rw.headers {
		for _, v := range vv {
			rw.ResponseWriter.Header().Add(k, v)
		}
	}
	rw.ResponseWriter.WriteHeader(rw.status)
	_, _ = rw.ResponseWriter.Write(rw.body)
}
