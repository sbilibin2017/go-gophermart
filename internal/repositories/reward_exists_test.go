package repositories

import (
	"context"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRewardExistsRepository_Exists_True(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := NewMockRewardExistsQuerier(ctrl)
	repo := NewRewardExistsRepository(mockQuerier)

	ctx := context.Background()
	rewardID := "reward-123"
	expectedExists := true

	argMap := map[string]any{
		"reward_id": rewardID,
	}

	mockQuerier.EXPECT().
		Query(ctx, gomock.AssignableToTypeOf(new(bool)), rewardExistsByIDQuery, argMap).
		DoAndReturn(func(ctx context.Context, dest any, query string, args map[string]any) error {
			ptr := dest.(*bool)
			*ptr = expectedExists
			return nil
		})

	exists, err := repo.Exists(ctx, rewardID)
	require.NoError(t, err)
	require.True(t, exists)
}

func TestRewardExistsRepository_Exists_False(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := NewMockRewardExistsQuerier(ctrl)
	repo := NewRewardExistsRepository(mockQuerier)

	ctx := context.Background()
	rewardID := "nonexistent-reward"
	expectedExists := false

	argMap := map[string]any{
		"reward_id": rewardID,
	}

	mockQuerier.EXPECT().
		Query(ctx, gomock.AssignableToTypeOf(new(bool)), rewardExistsByIDQuery, argMap).
		DoAndReturn(func(ctx context.Context, dest any, query string, args map[string]any) error {
			ptr := dest.(*bool)
			*ptr = expectedExists
			return nil
		})

	exists, err := repo.Exists(ctx, rewardID)
	require.NoError(t, err)
	require.False(t, exists)
}

func TestRewardExistsRepository_Exists_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := NewMockRewardExistsQuerier(ctrl)
	repo := NewRewardExistsRepository(mockQuerier)

	ctx := context.Background()
	rewardID := "reward-err"
	expectedErr := errors.New("query failed")

	argMap := map[string]any{
		"reward_id": rewardID,
	}

	mockQuerier.EXPECT().
		Query(ctx, gomock.Any(), rewardExistsByIDQuery, argMap).
		Return(expectedErr)

	exists, err := repo.Exists(ctx, rewardID)
	require.Error(t, err)
	require.EqualError(t, err, expectedErr.Error())
	require.False(t, exists)
}
