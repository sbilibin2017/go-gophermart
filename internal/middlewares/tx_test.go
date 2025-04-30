package middlewares_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/stretchr/testify/assert"
)

type contextTestKey string

func dummyTxSetter(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, contextTestKey("tx"), tx)
}

func TestTxMiddleware_CommitsOnSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mock.ExpectBegin()
	mock.ExpectCommit()

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	middleware := middlewares.TxMiddleware(sqlxDB, dummyTxSetter)
	middleware(handler).ServeHTTP(rec, req)

	assert.True(t, handlerCalled)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTxMiddleware_RollbackOnPanic(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mock.ExpectBegin()
	mock.ExpectRollback()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("forced panic")
	})

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	middleware := middlewares.TxMiddleware(sqlxDB, dummyTxSetter)

	assert.Panics(t, func() {
		middleware(handler).ServeHTTP(rec, req)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTxMiddleware_BeginTxFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called on failed BeginTxx")
	})

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	middleware := middlewares.TxMiddleware(sqlxDB, dummyTxSetter)
	middleware(handler).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}
