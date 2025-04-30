package middlewares

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/db"
)

func TxMiddleware(
	conn *sqlx.DB,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := db.WithTx(conn, func(tx *sqlx.Tx) error {
				ctx := contextutils.SetTx(r.Context(), tx)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return nil
			})
			if err != nil {
				return
			}
		})
	}
}
