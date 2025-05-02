package middlewares

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestGetTxFromContext_Success(t *testing.T) {
	ctx := context.Background()
	tx := &sqlx.Tx{}
	ctx = setTx(ctx, tx)
	result, err := GetTxFromContext(ctx)
	assert.NoError(t, err)
	assert.Equal(t, tx, result)
}

func TestGetTxFromContext_Error(t *testing.T) {
	ctx := context.Background()
	result, err := GetTxFromContext(ctx)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestSetTx(t *testing.T) {
	ctx := context.Background()
	tx := &sqlx.Tx{}
	ctx = setTx(ctx, tx)
	result, err := GetTxFromContext(ctx)
	assert.NoError(t, err)
	assert.Equal(t, tx, result)
}

func TestTxMiddleware_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectBegin()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		GetTxFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	middleware := TxMiddleware(sqlxDB)(handler)
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTxMiddleware_BeginTxxError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectBegin().WillReturnError(errors.New("db error"))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		GetTxFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	middleware := TxMiddleware(sqlxDB)(handler)
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}
