package storage

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContextWithTx(t *testing.T) {
	ctx := context.Background()
	mockTx := &sqlx.Tx{}
	ctx = ContextWithTx(ctx, mockTx)
	tx, ok := TxFromContext(ctx)
	assert.True(t, ok, "Expected to find a transaction in context")
	assert.Equal(t, mockTx, tx, "Expected to retrieve the same transaction from context")
}

func TestWithTx_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	mock.ExpectBegin()
	mock.ExpectCommit()
	err = WithTx(context.Background(), sqlxDB, func(tx *sqlx.Tx) error {
		assert.NotNil(t, tx, "Expected a non-nil transaction")
		return nil
	})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWithTx_RollbackOnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	mock.ExpectBegin()
	mock.ExpectRollback()
	err = WithTx(context.Background(), sqlxDB, func(tx *sqlx.Tx) error {
		assert.NotNil(t, tx, "Expected a non-nil transaction")
		return assert.AnError
	})
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWithTx_NilDB(t *testing.T) {
	err := WithTx(context.Background(), nil, func(tx *sqlx.Tx) error {
		t.Fatal("Expected WithTx to return early and not enter the function")
		return nil
	})
	assert.NoError(t, err)
}

func TestWithTx_BeginTxx_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	mock.ExpectBegin().WillReturnError(assert.AnError)
	err = WithTx(context.Background(), sqlxDB, func(tx *sqlx.Tx) error {
		t.Fatal("Expected WithTx to return early due to BeginTxx error")
		return nil
	})
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
