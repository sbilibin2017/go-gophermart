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

func TestOrderSaveRepository_Save(t *testing.T) {
	tests := []struct {
		name          string
		order         *OrderSaveDB
		mockExpect    func(mock sqlmock.Sqlmock, order *OrderSaveDB)
		expectedError bool
	}{
		{
			name: "Success",
			order: &OrderSaveDB{
				Number:  12345,
				Status:  StatusNew,
				Accrual: 100.50,
			},
			mockExpect: func(mock sqlmock.Sqlmock, order *OrderSaveDB) {
				mock.ExpectExec(`INSERT INTO orders \(number, status, accrual, created_at, updated_at\)`).
					WithArgs(order.Number, order.Status, order.Accrual).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name: "Error - Database Error",
			order: &OrderSaveDB{
				Number:  12345,
				Status:  StatusNew,
				Accrual: 100.50,
			},
			mockExpect: func(mock sqlmock.Sqlmock, order *OrderSaveDB) {
				mock.ExpectExec(`INSERT INTO orders \(number, status, accrual, created_at, updated_at\)`).
					WithArgs(order.Number, order.Status, order.Accrual).
					WillReturnError(fmt.Errorf("db error"))
			},
			expectedError: true,
		},
		{
			name: "Error - Invalid Status",
			order: &OrderSaveDB{
				Number:  12345,
				Status:  "INVALID_STATUS", // Invalid status
				Accrual: 100.50,
			},
			mockExpect: func(mock sqlmock.Sqlmock, order *OrderSaveDB) {
				mock.ExpectExec(`INSERT INTO orders \(number, status, accrual, created_at, updated_at\)`).
					WithArgs(order.Number, order.Status, order.Accrual).
					WillReturnError(fmt.Errorf("db error"))
			},
			expectedError: true,
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
			repo := NewOrderSaveRepository(sqlxDB)

			tt.mockExpect(mock, tt.order)

			err = repo.Save(context.Background(), tt.order)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
