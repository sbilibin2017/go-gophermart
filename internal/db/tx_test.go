package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestTx_Do(t *testing.T) {
	type testCase struct {
		name         string
		setupMock    func(mock sqlmock.Sqlmock)
		isErr        bool
		expectCalled bool
		fn           func(tx *sqlx.Tx) error
	}

	tests := []testCase{
		{
			name: "успешное выполнение транзакции",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectCommit()
			},
			fn: func(tx *sqlx.Tx) error {
				// Можно сюда добавить какие-то проверки или работу с tx, если нужно
				return nil
			},
			isErr:        false,
			expectCalled: true,
		},
		{
			name: "ошибка при начале транзакции",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(fmt.Errorf("begin error"))
			},
			fn: func(tx *sqlx.Tx) error {
				return nil
			},
			isErr:        true,
			expectCalled: false,
		},
		{
			name: "ошибка внутри функции fn",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectRollback()
			},
			fn: func(tx *sqlx.Tx) error {
				return fmt.Errorf("some error")
			},
			isErr:        true,
			expectCalled: true,
		},
		{
			name: "ошибка при коммите транзакции",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectCommit().WillReturnError(fmt.Errorf("commit failed"))
			},
			fn: func(tx *sqlx.Tx) error {
				return nil
			},
			isErr:        true,
			expectCalled: true,
		},
		{
			name:         "отсутствие подключения к бд",
			setupMock:    nil, // не нужен mock
			fn:           func(tx *sqlx.Tx) error { return nil },
			isErr:        true,
			expectCalled: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var (
				db   *sqlx.DB
				mock sqlmock.Sqlmock
				err  error
			)

			if tc.name != "отсутствие подключения к бд" {
				rawDB, sqlmock, mockErr := sqlmock.New()
				assert.NoError(t, mockErr)
				defer rawDB.Close()

				db = sqlx.NewDb(rawDB, "sqlmock")
				mock = sqlmock
			} else {
				db = nil
			}

			tx := &Tx{db: db}

			if tc.setupMock != nil {
				tc.setupMock(mock)
			}

			called := false
			err = tx.Do(context.Background(), func(tx *sqlx.Tx) error {
				called = true
				return tc.fn(tx)
			})

			if tc.isErr {
				assert.Error(t, err)
				if tc.name == "отсутствие подключения к бд" {
					assert.Contains(t, err.Error(), "отсутствие подключения к бд")
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectCalled, called, "ожидание вызова fn не совпало")

			if mock != nil {
				assert.NoError(t, mock.ExpectationsWereMet(), "все ожидаемые SQL действия должны быть выполнены")
			}
		})
	}
}
