package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockDatabaseURIGetter struct {
	URI string
}

func (m *MockDatabaseURIGetter) GetDatabaseURI() string {
	return m.URI
}

func TestNewDB_Success(t *testing.T) {
	_, host, port, cleanup := Setup(t)
	defer cleanup()
	mockGetter := &MockDatabaseURIGetter{
		URI: "postgres://testuser:testpassword@" + host + ":" + port + "/testdb?sslmode=disable",
	}
	t.Log("Mock Database URI: ", mockGetter.URI)
	db, err := NewDB(mockGetter)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	err = db.Ping()
	assert.NoError(t, err)
}

func TestNewDB_Error(t *testing.T) {
	_, _, _, cleanup := Setup(t)
	defer cleanup()
	mockGetter := &MockDatabaseURIGetter{URI: "invalid-uri"}
	db, err := NewDB(mockGetter)
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestNewDB_InvalidURI(t *testing.T) {
	_, _, _, cleanup := Setup(t)
	defer cleanup()
	mockGetter := &MockDatabaseURIGetter{URI: "postgres://invalid-uri"}
	db, err := NewDB(mockGetter)
	assert.Error(t, err)
	assert.Nil(t, db)
}
