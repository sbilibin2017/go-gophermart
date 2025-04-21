package middlewares

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
)

func TxMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Starting a new transaction")

			tx, err := db.Begin()
			if err != nil {
				log.Printf("Failed to start transaction: %v", err)
				http.Error(w, "could not start transaction", http.StatusInternalServerError)
				return
			}

			ctx := contextutils.SetTx(r.Context(), tx)
			rw := newBufferedResponseWriter(w)

			next.ServeHTTP(rw, r.WithContext(ctx))

			log.Printf("Request completed with status code: %d", rw.status)

			if rw.status >= 400 {
				log.Println("Rolling back transaction due to client or server error")
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					log.Printf("Failed to rollback transaction: %v", rollbackErr)
					http.Error(w, "could not rollback transaction", http.StatusInternalServerError)
					return
				}
				log.Println("Transaction rolled back successfully")
				rw.FlushToUnderlying() // всё равно отдаём клиенту то, что хотел хендлер
			} else {
				log.Println("Committing transaction")
				if commitErr := tx.Commit(); commitErr != nil {
					log.Printf("Failed to commit transaction: %v", commitErr)
					http.Error(w, "could not commit transaction", http.StatusInternalServerError)
					return
				}
				log.Println("Transaction committed successfully")
				rw.FlushToUnderlying()
			}
		})
	}
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
