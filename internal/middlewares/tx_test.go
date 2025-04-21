package middlewares

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func createTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open SQLite database: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT);`)
	if err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}
	return db
}

func TestTxMiddleware(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()
	txMiddleware := TxMiddleware(db)
	handler := txMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx, ok := contextutils.GetTx(r.Context())
		if !ok || tx == nil {
			http.Error(w, "Transaction not found in context", http.StatusInternalServerError)
			return
		}
		_, err := tx.Exec("INSERT INTO users (name) VALUES (?)", "John Doe")
		if err != nil {
			http.Error(w, fmt.Sprintf("error inserting user: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User inserted"))
	}))
	ts := httptest.NewServer(handler)
	defer ts.Close()
	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("failed to make a GET request: %v", err)
	}
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE name = ?", "John Doe").Scan(&count)
	if err != nil {
		t.Fatalf("failed to query the database: %v", err)
	}
	assert.Equal(t, 1, count, "Expected one user in the database")
}

func TestTxMiddleware_Rollback(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()
	txMiddleware := TxMiddleware(db)
	handler := txMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx, ok := contextutils.GetTx(r.Context())
		if !ok || tx == nil {
			http.Error(w, "Transaction not found in context", http.StatusInternalServerError)
			return
		}
		_, err := tx.Exec("INSERT INTO users (wrong_column) VALUES (?)", "Jane Doe")
		if err != nil {
			http.Error(w, fmt.Sprintf("error inserting user: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User inserted"))
	}))
	ts := httptest.NewServer(handler)
	defer ts.Close()
	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("failed to make a GET request: %v", err)
	}
	defer res.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE name = ?", "Jane Doe").Scan(&count)
	if err != nil {
		t.Fatalf("failed to query the database: %v", err)
	}
	assert.Equal(t, 0, count, "Expected no users to be inserted after rollback")
}

func TestTxMiddleware_RollbackError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()
	mock.ExpectBegin()
	mock.ExpectRollback().WillReturnError(fmt.Errorf("mock rollback failure"))
	txMiddleware := TxMiddleware(db)
	handler := txMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Simulated error", http.StatusInternalServerError)
	}))
	ts := httptest.NewServer(handler)
	defer ts.Close()
	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("failed to make a GET request: %v", err)
	}
	defer res.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unmet expectations: %v", err)
	}
}

func TestTxMiddleware_CommitError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectBegin()
	mock.ExpectCommit().WillReturnError(fmt.Errorf("mock commit failure"))
	middleware := TxMiddleware(db)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx, _ := contextutils.GetTx(r.Context())
		assert.NotNil(t, tx)
		_, _ = w.Write([]byte("OK"))
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTxMiddleware_BeginError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	mock.ExpectBegin().WillReturnError(errors.New("mock begin failure"))
	middleware := TxMiddleware(db)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called on begin error")
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.Contains(t, rec.Body.String(), "could not start transaction")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestBufferedResponseWriter_WriteHeader(t *testing.T) {
	rr := httptest.NewRecorder()
	rw := &bufferedResponseWriter{ResponseWriter: rr}
	rw.WriteHeader(http.StatusTeapot)
	require.True(t, rw.wrote)
	require.Equal(t, http.StatusTeapot, rw.status)
	rw.WriteHeader(http.StatusInternalServerError)
	require.Equal(t, http.StatusTeapot, rw.status, "status should not change after first WriteHeader call")
}
