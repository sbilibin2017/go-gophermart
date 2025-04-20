package db

import (
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDB_Success(t *testing.T) {
	originalFactory := connFactory
	t.Cleanup(func() {
		connFactory = originalFactory
	})
	connFactory = func(driverName, dsn string) (*sqlx.DB, error) {
		return sqlx.NewDb(nil, driverName), nil
	}
	dbConn, err := NewDB("fake-dsn")
	require.NoError(t, err)
	assert.NotNil(t, dbConn)
}

func TestNewDB_Error(t *testing.T) {
	originalFactory := connFactory
	t.Cleanup(func() {
		connFactory = originalFactory
	})
	expectedErr := errors.New("connection failed")
	connFactory = func(driverName, dsn string) (*sqlx.DB, error) {
		return nil, expectedErr
	}
	dbConn, err := NewDB("invalid-dsn")
	require.ErrorIs(t, err, expectedErr)
	assert.Nil(t, dbConn)
}
