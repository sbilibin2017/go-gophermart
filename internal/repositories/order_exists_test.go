package repositories

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func init() {
	logger.Init(zapcore.InfoLevel)
}

func TestOrderExistRepository_Exists(t *testing.T) {
	tests := []struct {
		name           string
		filter         *OrderExistsFilterDB
		mockExpect     func(mock sqlmock.Sqlmock, filter *OrderExistsFilterDB)
		expectedExists bool
		expectedError  bool
	}{
		{
			name:   "Success",
			filter: &OrderExistsFilterDB{Number: 12345},
			mockExpect: func(mock sqlmock.Sqlmock, filter *OrderExistsFilterDB) {
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM orders WHERE number = \$1\)`).
					WithArgs(filter.Number).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
			},
			expectedExists: true,
			expectedError:  false,
		},
		{
			name:   "Error",
			filter: &OrderExistsFilterDB{Number: 12345},
			mockExpect: func(mock sqlmock.Sqlmock, filter *OrderExistsFilterDB) {
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM orders WHERE number = \$1\)`).
					WithArgs(filter.Number).
					WillReturnError(fmt.Errorf("db error"))
			},
			expectedExists: false,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to open mock database: %v", err)
			}
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "pgx")
			repo := NewOrderExistRepository(sqlxDB)

			tt.mockExpect(mock, tt.filter)

			exists, err := repo.Exists(context.Background(), tt.filter)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedExists, exists)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
