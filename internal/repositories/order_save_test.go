package repositories

import (
	"context"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestOrderSaveRepository_SaveSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExecutor := NewMockOrderExecutor(ctrl)
	repo := NewOrderSaveRepository(mockExecutor)

	ctx := context.Background()
	orderID := "12345"
	status := "PROCESSED"
	accrual := 100.5

	argMap := map[string]any{
		"order_id": orderID,
		"status":   status,
		"accrual":  accrual,
	}

	mockExecutor.EXPECT().
		Execute(ctx, orderSaveQuery, argMap).
		Return(nil)

	err := repo.Save(ctx, orderID, status, accrual)
	require.NoError(t, err)
}

func TestOrderSaveRepository_SaveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExecutor := NewMockOrderExecutor(ctrl)
	repo := NewOrderSaveRepository(mockExecutor)

	ctx := context.Background()
	orderID := "12345"
	status := "PROCESSED"
	accrual := 100.5

	expectedErr := errors.New("db error")

	argMap := map[string]any{
		"order_id": orderID,
		"status":   status,
		"accrual":  accrual,
	}

	mockExecutor.EXPECT().
		Execute(ctx, orderSaveQuery, argMap).
		Return(expectedErr)

	err := repo.Save(ctx, orderID, status, accrual)
	require.EqualError(t, err, expectedErr.Error())
}
