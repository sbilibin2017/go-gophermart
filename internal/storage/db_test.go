package storage

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDB_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectPing()

	connOpener = func(driverName, dsn string) (*sql.DB, error) {
		return db, nil
	}

	result, err := NewDB("mock_dsn")
	require.NoError(t, err)
	assert.Equal(t, db, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestNewDB_OpenError(t *testing.T) {
	expectedErr := errors.New("connection failed")

	connOpener = func(driverName, dsn string) (*sql.DB, error) {
		return nil, expectedErr
	}

	db, err := NewDB("mock_dsn")
	assert.Nil(t, db)
	assert.ErrorIs(t, err, expectedErr)
}
