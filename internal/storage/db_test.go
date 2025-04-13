package storage

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
	_, host, port, cleanup := SetupPostgresContainer(t)
	defer cleanup()
	mockGetter := &MockDatabaseURIGetter{
		URI: "postgres://testuser:testpassword@" + host + ":" + port + "/testdb?sslmode=disable",
	}
	t.Log("Mock Database URI: ", mockGetter.URI)
	newDb, err := NewDB(mockGetter)
	assert.NoError(t, err)
	assert.NotNil(t, newDb)
	err = newDb.Ping()
	assert.NoError(t, err)
}

func TestNewDB_Error(t *testing.T) {
	_, _, _, cleanup := SetupPostgresContainer(t)
	defer cleanup()
	mockGetter := &MockDatabaseURIGetter{URI: "invalid-uri"}
	newDb, err := NewDB(mockGetter)
	assert.Error(t, err)
	assert.Nil(t, newDb)
}

func TestNewDB_InvalidURI(t *testing.T) {
	_, _, _, cleanup := SetupPostgresContainer(t)
	defer cleanup()
	mockGetter := &MockDatabaseURIGetter{URI: "postgres://invalid-uri"}
	newDb, err := NewDB(mockGetter)
	assert.Error(t, err)
	assert.Nil(t, newDb)
}
