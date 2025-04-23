package repositories

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/types"
	"github.com/stretchr/testify/require"
)

func TestOrderGetByIDRepository_GetByIDSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := NewMockOrderGetByIDQuerier(ctrl)
	repo := NewOrderGetByIDRepository(mockQuerier)

	ctx := context.Background()
	orderID := "12345"
	fields := []string{"order_id", "status", "accrual", "created_at", "updated_at"}
	var a int64 = 10
	expectedQuery := "SELECT order_id, status, accrual, created_at, updated_at FROM orders WHERE order_id = :order_id"
	expectedArgs := map[string]any{"order_id": orderID}
	expectedResult := types.OrderDB{
		OrderID:   orderID,
		Status:    "PROCESSED",
		Accrual:   &a,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 👇 Правильный SetArg и проверка указателя
	mockQuerier.EXPECT().
		Query(gomock.Any(), gomock.AssignableToTypeOf(&types.OrderDB{}), expectedQuery, expectedArgs).
		DoAndReturn(func(ctx context.Context, dest any, query string, args map[string]any) error {
			ptr := dest.(*types.OrderDB)
			*ptr = expectedResult
			return nil
		})

	result, err := repo.GetByID(ctx, orderID, fields)
	require.NoError(t, err)
	require.Equal(t, &expectedResult, result)
}

func TestOrderGetByIDRepository_GetByIDQueryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := NewMockOrderGetByIDQuerier(ctrl)
	repo := NewOrderGetByIDRepository(mockQuerier)

	ctx := context.Background()
	orderID := "12345"
	fields := []string{"order_id", "status"}

	expectedQuery := "SELECT order_id, status FROM orders WHERE order_id = :order_id"
	expectedArgs := map[string]any{"order_id": orderID}
	expectedErr := errors.New("query execution failed")

	mockQuerier.EXPECT().
		Query(gomock.Any(), gomock.AssignableToTypeOf(&types.OrderDB{}), expectedQuery, expectedArgs).
		Return(expectedErr)

	result, err := repo.GetByID(ctx, orderID, fields)
	require.Nil(t, result)
	require.EqualError(t, err, expectedErr.Error())
}
