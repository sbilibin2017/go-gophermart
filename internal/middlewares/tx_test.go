package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestTxMiddleware_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "postgres")

	mock.ExpectBegin()
	mock.ExpectCommit()

	txSetter := func(ctx context.Context, tx *sqlx.Tx) context.Context {
		return SetTx(ctx, tx)
	}

	middleware := TxMiddleware(sqlxDB, txSetter)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mwHandler := middleware(handler)

	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := &mockResponseWriter{}

	mwHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTxMiddleware_RollbackOnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "postgres")

	mock.ExpectBegin()
	mock.ExpectRollback()

	txSetter := func(ctx context.Context, tx *sqlx.Tx) context.Context {
		return SetTx(ctx, tx)
	}

	middleware := TxMiddleware(sqlxDB, txSetter)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	})

	mwHandler := middleware(handler)

	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := &mockResponseWriter{
		ResponseWriter: httptest.NewRecorder(),
	}

	mwHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

type mockResponseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *mockResponseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
}

func (rw *mockResponseWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func TestTxMiddleware_StartTransaction_Failure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "postgres")

	mock.ExpectBegin().WillReturnError(fmt.Errorf("failed to start transaction"))

	txSetter := func(ctx context.Context, tx *sqlx.Tx) context.Context {
		return SetTx(ctx, tx)
	}

	middleware := TxMiddleware(sqlxDB, txSetter)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called if transaction fails")
	})

	mwHandler := middleware(handler)

	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := &mockResponseWriter{
		ResponseWriter: httptest.NewRecorder(),
	}

	mwHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTx_Success(t *testing.T) {
	mockTx := &sqlx.Tx{}
	ctx := SetTx(context.Background(), mockTx)

	tx, ok := GetTx(ctx)

	assert.Equal(t, mockTx, tx)
	assert.True(t, ok)
}

func TestGetTx_NoTxInContext(t *testing.T) {
	ctx := context.Background()

	tx, ok := GetTx(ctx)

	assert.Nil(t, tx)
	assert.False(t, ok)
}

func TestTxMiddleware_RollbackFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "mysql")

	mock.ExpectBegin()

	mock.ExpectRollback().WillReturnError(fmt.Errorf("rollback failed"))

	txSetter := func(ctx context.Context, tx *sqlx.Tx) context.Context {
		return SetTx(ctx, tx)
	}

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/test", nil)
	w := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Server error", http.StatusInternalServerError)
	})

	txMiddleware := TxMiddleware(sqlxDB, txSetter)
	handler := txMiddleware(nextHandler)

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestTxMiddleware_CommitFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "mysql")

	mock.ExpectBegin()

	mock.ExpectCommit().WillReturnError(fmt.Errorf("commit failed"))

	txSetter := func(ctx context.Context, tx *sqlx.Tx) context.Context {
		return SetTx(ctx, tx)
	}

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/test", nil)
	w := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	txMiddleware := TxMiddleware(sqlxDB, txSetter)
	handler := txMiddleware(nextHandler)

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestBufferedResponseWriter_WriteHeader(t *testing.T) {
	recorder := httptest.NewRecorder()
	rw := newBufferedResponseWriter(recorder)

	rw.WriteHeader(http.StatusOK)
	assert.Equal(t, http.StatusOK, rw.status)

	rw.WriteHeader(http.StatusBadRequest)
	assert.Equal(t, http.StatusOK, rw.status)
}

func TestBufferedResponseWriter_Write(t *testing.T) {
	recorder := httptest.NewRecorder()
	rw := newBufferedResponseWriter(recorder)

	n, err := rw.Write([]byte("Test body"))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rw.status)
	assert.Equal(t, 9, n)

	assert.Equal(t, []byte("Test body"), rw.body)

	n, err = rw.Write([]byte("More body"))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rw.status)
	assert.Equal(t, 9, n)

	expectedBody := []byte("Test bodyMore body")
	assert.Equal(t, expectedBody, rw.body)

	totalBytesWritten := len(rw.body)
	assert.Equal(t, 18, totalBytesWritten)
}
