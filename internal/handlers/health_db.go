package handlers

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func HealthDBHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := db.Ping()
		if err != nil {
			log.Printf("Database connection is unhealthy: %v\n", err)
			http.Error(w, "Database connection is unhealthy", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Database connection is healthy"))
	}
}
