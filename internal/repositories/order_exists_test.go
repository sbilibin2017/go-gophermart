package repositories

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderExistsRepository_ExistsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := NewMockOrderExistsQuerier(ctrl)
	orderID := "12345"
	expectedQuery := orderExistsByIDQuery
	argMap := map[string]any{
		"order_id": orderID,
	}
	mockQuerier.EXPECT().
		Query(gomock.Any(), gomock.Any(), expectedQuery, argMap).
		SetArg(1, true).
		Return(nil)

	repo := NewOrderExistsRepository(mockQuerier)
	exists, err := repo.Exists(context.Background(), orderID)

	require.NoError(t, err)
	assert.True(t, exists)
}

func TestOrderExistsRepository_ExistsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := NewMockOrderExistsQuerier(ctrl)
	orderID := "12345"
	expectedQuery := orderExistsByIDQuery
	argMap := map[string]any{
		"order_id": orderID,
	}
	mockQuerier.EXPECT().
		Query(gomock.Any(), gomock.Any(), expectedQuery, argMap).
		Return(assert.AnError)

	repo := NewOrderExistsRepository(mockQuerier)
	exists, err := repo.Exists(context.Background(), orderID)

	assert.Error(t, err)
	assert.False(t, exists)
}

func TestOrderExistsRepository_ExistsNoRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := NewMockOrderExistsQuerier(ctrl)
	orderID := "12345"
	expectedQuery := orderExistsByIDQuery
	argMap := map[string]any{
		"order_id": orderID,
	}
	mockQuerier.EXPECT().
		Query(gomock.Any(), gomock.Any(), expectedQuery, argMap).
		SetArg(1, false).
		Return(nil)

	repo := NewOrderExistsRepository(mockQuerier)
	exists, err := repo.Exists(context.Background(), orderID)

	require.NoError(t, err)
	assert.False(t, exists)
}
