package middlewares

import (
	"database/sql"
	"log"
	"net/http"
)

func TxMiddleware(
	db *sql.DB,
	txFactory func(db *sql.DB, op func(tx *sql.Tx) error) error, // txFactory теперь ожидает op как аргумент
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := txFactory(db, func(tx *sql.Tx) error {
				return nil
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Println("Error in transaction: ", err)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
