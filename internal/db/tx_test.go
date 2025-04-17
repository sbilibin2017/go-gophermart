package db

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestWithTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "postgres")
	ctx := context.Background()

	tests := []struct {
		name        string
		mockSetup   func()
		txFunc      func(tx *sqlx.Tx) error
		expectedErr error
		db          *sqlx.DB
	}{
		{
			name: "success",
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectCommit()
			},
			txFunc: func(tx *sqlx.Tx) error {
				return nil
			},
			expectedErr: nil,
			db:          sqlxDB,
		},
		{
			name: "rollback_on_error",
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectRollback()
			},
			txFunc: func(tx *sqlx.Tx) error {
				return assert.AnError
			},
			expectedErr: assert.AnError,
			db:          sqlxDB,
		},
		{
			name: "commit_on_success",
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectCommit()
			},
			txFunc: func(tx *sqlx.Tx) error {
				return nil
			},
			expectedErr: nil,
			db:          sqlxDB,
		},
		{
			name: "nil_db",
			mockSetup: func() {
			},
			txFunc: func(tx *sqlx.Tx) error {
				return nil
			},
			expectedErr: nil,
			db:          nil,
		},
		{
			name: "error_begin_tx",
			mockSetup: func() {
				mock.ExpectBegin().WillReturnError(assert.AnError)
			},
			txFunc: func(tx *sqlx.Tx) error {
				return nil
			},
			expectedErr: assert.AnError,
			db:          sqlxDB,
		},
		{
			name: "error_commit_tx",
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectCommit().WillReturnError(assert.AnError)
			},
			txFunc: func(tx *sqlx.Tx) error {
				return nil
			},
			expectedErr: assert.AnError,
			db:          sqlxDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.db != nil {
				tt.mockSetup()
			}
			err := WithTx(ctx, tt.db, tt.txFunc)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
