package db

import (
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestNewDB_TableDriven(t *testing.T) {
	logger.Init()

	type testCase struct {
		name        string
		mockOpener  func(driverName, dsn string) (*sqlx.DB, error)
		expectedDB  *sqlx.DB
		expectedErr error
	}

	fakeDB := &sqlx.DB{}
	tests := []testCase{
		{
			name: "success",
			mockOpener: func(driverName, dsn string) (*sqlx.DB, error) {
				return fakeDB, nil
			},
			expectedDB:  fakeDB,
			expectedErr: nil,
		},
		{
			name: "connection error",
			mockOpener: func(driverName, dsn string) (*sqlx.DB, error) {
				return nil, errors.New("connection error")
			},
			expectedDB:  nil,
			expectedErr: errors.New("connection error"),
		},
	}

	originalOpener := opener
	defer func() { opener = originalOpener }()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opener = tc.mockOpener

			db, err := NewDB("test-dsn")

			assert.Equal(t, tc.expectedDB, db)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
