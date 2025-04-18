package storage

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithTrx_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO users").WithArgs("John Doe").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err = WithTx(db, func(tx *sql.Tx) error {
		_, err := tx.Exec("INSERT INTO users (name) VALUES ($1)", "John Doe")
		return err
	})
	require.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestWithTx_BeginError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	mock.ExpectBegin().WillReturnError(errors.New("begin error"))
	err = WithTx(db, func(tx *sql.Tx) error {
		return nil
	})
	assert.Error(t, err)
	assert.EqualError(t, err, "begin error")
}

func TestWithTx_OpError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	mock.ExpectBegin()
	mock.ExpectRollback()
	err = WithTx(db, func(tx *sql.Tx) error {
		return errors.New("operation error")
	})
	assert.Error(t, err)
	assert.EqualError(t, err, "operation error")
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestWithTx_RollbackError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	mock.ExpectBegin()
	mock.ExpectRollback().WillReturnError(errors.New("rollback error"))
	err = WithTx(db, func(tx *sql.Tx) error {
		return errors.New("operation error")
	})
	assert.Error(t, err)
	assert.EqualError(t, err, "rollback error")
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestWithTx_DBIsNil(t *testing.T) {
	err := WithTx(nil, func(tx *sql.Tx) error {
		t.Fatal("This should not be reached")
		return nil
	})
	assert.NoError(t, err)
}
