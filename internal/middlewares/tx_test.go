package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestTxMiddleware_Success(t *testing.T) {
	logger.Init()
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
	logger.Init()
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
	logger.Init()
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
	logger.Init()
	mockTx := &sqlx.Tx{}
	ctx := SetTx(context.Background(), mockTx)

	tx, ok := GetTx(ctx)

	assert.Equal(t, mockTx, tx)
	assert.True(t, ok)
}

func TestGetTx_NoTxInContext(t *testing.T) {
	logger.Init()
	ctx := context.Background()

	tx, ok := GetTx(ctx)

	assert.Nil(t, tx)
	assert.False(t, ok)
}
