package storage

import (
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func fakeConnectSuccess(driverName, dsn string) (*sqlx.DB, error) {
	return &sqlx.DB{}, nil
}

func fakeConnectError(driverName, dsn string) (*sqlx.DB, error) {
	return nil, errors.New("mocked connection error")
}

func TestNewDB_Success(t *testing.T) {
	originalConnect := connect
	connect = fakeConnectSuccess
	defer func() { connect = originalConnect }()
	db, err := NewDB("mock-dsn")
	assert.NoError(t, err)
	assert.NotNil(t, db)
}

func TestNewDB_Error(t *testing.T) {
	originalConnect := connect
	connect = fakeConnectError
	defer func() { connect = originalConnect }()
	db, err := NewDB("mock-dsn")
	assert.Error(t, err)
	assert.Nil(t, db)
}
