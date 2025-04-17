package db

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestWithTx_TableDriven(t *testing.T) {
	type testCase struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock)
		dbIsNil       bool
		fn            func(tx *sqlx.Tx) error
		expectedError error
		expectCalled  bool
	}

	tests := []testCase{
		{
			name: "success",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectCommit()
			},
			fn: func(tx *sqlx.Tx) error {
				return nil
			},
			expectedError: nil,
			expectCalled:  true,
		},
		{
			name: "fn returns error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectRollback()
			},
			fn: func(tx *sqlx.Tx) error {
				return errors.New("something went wrong")
			},
			expectedError: errors.New("something went wrong"),
			expectCalled:  true,
		},
		{
			name: "begin fails",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("begin failed"))
			},
			fn: func(tx *sqlx.Tx) error {
				return nil
			},
			expectedError: errors.New("begin failed"),
			expectCalled:  false,
		},
		{
			name:    "nil db",
			dbIsNil: true,
			fn: func(tx *sqlx.Tx) error {
				return nil
			},
			expectedError: nil,
			expectCalled:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var (
				db   *sqlx.DB
				mock sqlmock.Sqlmock
				err  error
			)

			if !tc.dbIsNil {
				var rawDB *sql.DB
				rawDB, mock, err = sqlmock.New()
				assert.NoError(t, err)
				defer rawDB.Close()

				db = sqlx.NewDb(rawDB, "sqlmock")

				if tc.setupMock != nil {
					tc.setupMock(mock)
				}
			}

			called := false
			err = WithTx(context.Background(), db, func(tx *sqlx.Tx) error {
				called = true
				return tc.fn(tx)
			})

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectCalled, called)

			if mock != nil {
				assert.NoError(t, mock.ExpectationsWereMet())
			}
		})
	}
}
