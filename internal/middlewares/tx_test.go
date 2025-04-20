package middlewares

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sbilibin2017/go-gophermart/internal/handlers/utils"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

func createTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT)`)
	require.NoError(t, err)
	return db
}

func TestTxMiddleware_Commit(t *testing.T) {
	logger.Logger, _ = zap.NewDevelopment()
	db := createTestDB(t)
	defer db.Close()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx := TxFromContext(r.Context())
		require.NotNil(t, tx)
		_, err := tx.ExecContext(r.Context(), "INSERT INTO users (name) VALUES (?)", "John")
		require.NoError(t, err)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
	middleware := TxMiddleware(db)(handler)
	req, err := http.NewRequest("GET", "/commit", nil)
	require.NoError(t, err)
	rec := httptest.NewRecorder()
	middleware.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	var name string
	row := db.QueryRow("SELECT name FROM users WHERE name = ?", "John")
	err = row.Scan(&name)
	require.NoError(t, err)
	assert.Equal(t, "John", name)
}

func TestTxMiddleware_Rollback(t *testing.T) {
	logger.Logger, _ = zap.NewDevelopment()
	db := createTestDB(t)
	defer db.Close()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx := TxFromContext(r.Context())
		require.NotNil(t, tx)
		_, err := tx.ExecContext(r.Context(), "INSERT INTO non_existing_table (name) VALUES (?)", "John")
		assert.Error(t, err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error"))
	})
	middleware := TxMiddleware(db)(handler)
	req, err := http.NewRequest("GET", "/rollback", nil)
	require.NoError(t, err)
	rec := httptest.NewRecorder()
	middleware.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM users")
	err = row.Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestTxMiddleware_CommitError(t *testing.T) {
	logger.Logger, _ = zap.NewDevelopment()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO users").
		WithArgs("John").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(errors.New("commit failed"))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx := TxFromContext(r.Context())
		require.NotNil(t, tx)
		_, err := tx.ExecContext(r.Context(), "INSERT INTO users (name) VALUES (?)", "John")
		require.NoError(t, err)
		if tx != nil {
			commitErr := tx.Commit()
			if commitErr != nil {
				logger.Logger.Error("Failed to commit transaction", zap.Error(commitErr))
				utils.ErrorInternalServerResponse(w, commitErr)
				return
			}
		}
	})
	middleware := TxMiddleware(db)(handler)
	req := httptest.NewRequest("POST", "/commit-error", nil)
	rec := httptest.NewRecorder()
	middleware.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBeginTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	req, err := http.NewRequest("GET", "http://example.com", nil)
	require.NoError(t, err)
	rw := &responseWriter{
		ResponseWriter: httptest.NewRecorder(),
	}
	mock.ExpectBegin()
	tx, err := beginTx(db, rw, req)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestBeginTx_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	mock.ExpectBegin().WillReturnError(sql.ErrConnDone)
	req, err := http.NewRequest("GET", "http://example.com", nil)
	require.NoError(t, err)
	rw := httptest.NewRecorder()
	tx, err := beginTx(db, rw, req)
	assert.Nil(t, tx)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.Equal(t, http.StatusInternalServerError, rw.Code)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestHandleTransactionEnd_TxNil(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	rw := httptest.NewRecorder()
	handleTransactionEnd(nil, http.StatusOK, rw)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

type MockResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *MockResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func TestHandleCommit_ErrorDuringCommit(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	mock.ExpectBegin()
	tx, err := db.Begin()
	require.NoError(t, err)
	mock.ExpectCommit().WillReturnError(errors.New("commit error"))
	rw := httptest.NewRecorder()
	handleCommit(tx, rw)
	assert.Equal(t, http.StatusInternalServerError, rw.Code)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestTxMiddleware_ErrorDuringTransactionBegin(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	mock.ExpectBegin().WillReturnError(errors.New("failed to begin transaction"))
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	middleware := TxMiddleware(db)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(rw, req)
	assert.Equal(t, http.StatusInternalServerError, rw.Code)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestTxFromContext(t *testing.T) {
	mockTx := &sql.Tx{}
	ctxWithTx := context.WithValue(context.Background(), txKey, mockTx)
	tx := TxFromContext(ctxWithTx)
	assert.Equal(t, mockTx, tx)
	ctxWithoutTx := context.Background()
	tx = TxFromContext(ctxWithoutTx)
	assert.Nil(t, tx)
}
