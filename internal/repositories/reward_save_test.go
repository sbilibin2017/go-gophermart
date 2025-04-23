package repositories

import (
	"context"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRewardSaveRepository_Save_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExecutor := NewMockRewardExecutor(ctrl)
	repo := NewRewardSaveRepository(mockExecutor)

	ctx := context.Background()
	rewardID := "reward-001"
	reward := int64(100)
	rewardType := "bonus"

	argMap := map[string]any{
		"reward_id":   rewardID,
		"reward":      reward,
		"reward_type": rewardType,
	}

	mockExecutor.EXPECT().
		Execute(ctx, rewardSaveQuery, argMap).
		Return(nil)

	err := repo.Save(ctx, rewardID, reward, rewardType)
	require.NoError(t, err)
}

func TestRewardSaveRepository_Save_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExecutor := NewMockRewardExecutor(ctrl)
	repo := NewRewardSaveRepository(mockExecutor)

	ctx := context.Background()
	rewardID := "reward-002"
	reward := int64(200)
	rewardType := "referral"

	argMap := map[string]any{
		"reward_id":   rewardID,
		"reward":      reward,
		"reward_type": rewardType,
	}

	expectedErr := errors.New("execute failed")

	mockExecutor.EXPECT().
		Execute(ctx, rewardSaveQuery, argMap).
		Return(expectedErr)

	err := repo.Save(ctx, rewardID, reward, rewardType)
	require.Error(t, err)
	require.EqualError(t, err, expectedErr.Error())
}
