package middlewares

import (
	"fmt"
	"net/http"
	"testing"

	"net/http/httptest"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTxMiddleware(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	txMiddleware := TxMiddleware(db)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx, ok := TxFromContext(r.Context())
		if !ok {
			http.Error(w, "no transaction found", http.StatusInternalServerError)
			return
		}

		if r.URL.Query().Get("fail") == "true" {
			_ = tx.Rollback()
			http.Error(w, "transaction rolled back", http.StatusInternalServerError)
			return
		}

		err := tx.Commit()
		if err != nil {
			http.Error(w, "failed to commit transaction", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	middlewareHandler := txMiddleware(handler)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectCommit()

		req, err := http.NewRequest(http.MethodGet, "http://localhost/?fail=false", nil)
		require.NoError(t, err)
		w := httptest.NewRecorder()

		middlewareHandler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Rollback", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectRollback()

		req, err := http.NewRequest(http.MethodGet, "http://localhost/?fail=true", nil)
		require.NoError(t, err)
		w := httptest.NewRecorder()

		middlewareHandler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "transaction rolled back")
	})

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestTxMiddleware_BeginError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin().WillReturnError(fmt.Errorf("failed to begin transaction"))

	txMiddleware := TxMiddleware(db)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "failed to begin transaction", http.StatusInternalServerError)
			return
		}
		defer tx.Commit()
	})

	middlewareHandler := txMiddleware(handler)

	req, err := http.NewRequest(http.MethodGet, "http://localhost/", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	middlewareHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to begin transaction")

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
