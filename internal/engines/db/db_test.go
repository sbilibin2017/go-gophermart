package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDatabaseURIGetter struct {
	mock.Mock
}

func (m *MockDatabaseURIGetter) GetDatabaseURI() string {
	args := m.Called()
	return args.String(0)
}

func TestNewDB_Success(t *testing.T) {
	mockURIGetter := new(MockDatabaseURIGetter)
	mockURIGetter.On("GetDatabaseURI").Return("postgres://postgres:password@localhost:5432/testdb?sslmode=disable")
	db := NewDB(mockURIGetter)
	assert.NotNil(t, db)

}
