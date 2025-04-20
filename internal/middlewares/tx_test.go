package middlewares_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/db"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestTxMiddleware(t *testing.T) {
	dbConn, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer dbConn.Close()

	_, err = dbConn.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT)")
	require.NoError(t, err)

	middleware := middlewares.TxMiddleware(dbConn)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx := db.TxFromContext(r.Context())
		require.NotNil(t, tx)

		_, err := tx.Exec("INSERT INTO users (name) VALUES (?)", "John Doe")
		require.NoError(t, err)

		w.WriteHeader(http.StatusOK)
	})

	rw := &mockResponseWriter{}
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	middleware(handler).ServeHTTP(rw, req)

	assert.Equal(t, http.StatusOK, rw.statusCode)

	var name string
	err = dbConn.Get(&name, "SELECT name FROM users WHERE id = 1")
	require.NoError(t, err)
	assert.Equal(t, "John Doe", name)
}

func TestTxMiddleware_RollbackOnError(t *testing.T) {
	dbConn, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer dbConn.Close()

	_, err = dbConn.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT)")
	require.NoError(t, err)

	middleware := middlewares.TxMiddleware(dbConn)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx := db.TxFromContext(r.Context())
		require.NotNil(t, tx)

		_, err := tx.Exec("INSERT INTO users (name) VALUES (?)", "John Doe")
		require.NoError(t, err)

		http.Error(w, "Simulated error", http.StatusBadRequest)
	})

	rw := &mockResponseWriter{}
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	middleware(handler).ServeHTTP(rw, req)

	assert.Equal(t, http.StatusBadRequest, rw.statusCode)

}

func TestTxMiddleware_BeginTransactionError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	mock.ExpectBegin().WillReturnError(fmt.Errorf("could not start transaction"))

	sqlxDB := sqlx.NewDb(dbConn, "postgres")

	middleware := middlewares.TxMiddleware(sqlxDB)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx, err := sqlxDB.BeginTxx(r.Context(), nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tx.Exec("INSERT INTO users (name) VALUES (?)", "John Doe")
		w.WriteHeader(http.StatusOK)
	})

	rw := &mockResponseWriter{}

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	middleware(handler).ServeHTTP(rw, req)

	assert.Equal(t, http.StatusInternalServerError, rw.statusCode)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

type mockResponseWriter struct {
	statusCode int
}

func (rw *mockResponseWriter) Header() http.Header {
	return http.Header{}
}

func (rw *mockResponseWriter) Write(data []byte) (int, error) {
	return len(data), nil
}

func (rw *mockResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
}
