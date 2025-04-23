package repositories

import (
	"context"
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/types"
	"github.com/stretchr/testify/require"
)

func TestRewardFilterILikeRepository_FilterILike_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := NewMockRewardFilterILikeQuerier(ctrl)
	repo := NewRewardFilterILikeRepository(mockQuerier)

	ctx := context.Background()
	rewardID := "bonus"
	fields := []string{"reward_id", "reward", "reward_type"}
	expectedQuery := "SELECT reward_id, reward, reward_type FROM rewards WHERE reward_id ILIKE :reward_id"
	expectedArgMap := map[string]any{"reward_id": "%bonus%"}

	expectedResult := &types.RewardDB{
		RewardID:   "123",
		Reward:     100,
		RewardType: "bonus",
	}

	mockQuerier.EXPECT().
		Query(ctx, gomock.AssignableToTypeOf(&expectedResult), expectedQuery, expectedArgMap).
		DoAndReturn(func(_ context.Context, dest any, _ string, _ map[string]any) error {
			ptr := dest.(**types.RewardDB)
			*ptr = expectedResult
			return nil
		})

	result, err := repo.FilterILike(ctx, rewardID, fields)
	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestRewardFilterILikeRepository_FilterILike_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := NewMockRewardFilterILikeQuerier(ctrl)
	repo := NewRewardFilterILikeRepository(mockQuerier)

	ctx := context.Background()
	rewardID := "fail"
	fields := []string{"reward_id"}
	expectedQuery := "SELECT reward_id FROM rewards WHERE reward_id ILIKE :reward_id"
	expectedArgMap := map[string]any{"reward_id": "%fail%"}

	expectedErr := fmt.Errorf("query failed")

	mockQuerier.EXPECT().
		Query(ctx, gomock.Any(), expectedQuery, expectedArgMap).
		Return(expectedErr)

	result, err := repo.FilterILike(ctx, rewardID, fields)
	require.Nil(t, result)
	require.EqualError(t, err, expectedErr.Error())
}
