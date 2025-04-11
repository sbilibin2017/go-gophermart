package handlers

import (
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Pinger interface {
	Ping() error
}

func GophermartPingHandler(p Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if p == nil {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
		if err := p.Ping(); err != nil {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Database connection successful"))
	}
}
