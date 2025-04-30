package routers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestRegisterAccrualRouter(t *testing.T) {
	db, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	mainRouter := chi.NewRouter()

	dummyHandler := func(body string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(body))
		})
	}

	mockMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}

	txMiddleware := func(db *sqlx.DB) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				tx, err := db.Beginx()
				if err != nil {
					http.Error(w, "Unable to begin transaction: "+err.Error(), http.StatusInternalServerError)
					return
				}
				defer tx.Rollback()

				ctx := context.WithValue(r.Context(), "tx", tx)
				next.ServeHTTP(w, r.WithContext(ctx))
				tx.Commit()
			})
		}
	}

	RegisterAccrualRouter(
		mainRouter,
		db,
		"/api",
		dummyHandler("order info"),
		dummyHandler("order created"),
		dummyHandler("good reward"),
		mockMiddleware,
		mockMiddleware,
		txMiddleware(db),
	)

	t.Run("GET /api/orders/{number}", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/orders/12345", nil)
		rec := httptest.NewRecorder()

		mainRouter.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "order info", rec.Body.String())
	})

	t.Run("POST /api/orders", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/orders", nil)
		rec := httptest.NewRecorder()

		mainRouter.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "order created", rec.Body.String())
	})

	t.Run("POST /api/goods", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/goods", nil)
		rec := httptest.NewRecorder()

		mainRouter.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "good reward", rec.Body.String())
	})
}
