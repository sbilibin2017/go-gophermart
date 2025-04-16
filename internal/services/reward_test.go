package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
	e "github.com/sbilibin2017/go-gophermart/internal/errors"
)

func TestRewardService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockExistsRepo := NewMockRewardExistsRepository(ctrl)
	mockSaveRepo := NewMockRewardSaveRepository(ctrl)

	rewardService := NewRewardService(mockExistsRepo, mockSaveRepo)

	testReward := &domain.Reward{
		Match:      "test-match",
		Reward:     100,
		RewardType: domain.Points,
	}

	t.Run("success - new reward", func(t *testing.T) {
		mockExistsRepo.
			EXPECT().
			Exists(ctx, testReward.Match).
			Return(false, nil)

		mockSaveRepo.
			EXPECT().
			Save(ctx, gomock.Any()).
			Return(nil)

		err := rewardService.Register(ctx, testReward)
		assert.NoError(t, err)
	})

	t.Run("error - reward already exists", func(t *testing.T) {
		mockExistsRepo.
			EXPECT().
			Exists(ctx, testReward.Match).
			Return(true, nil)

		err := rewardService.Register(ctx, testReward)
		assert.Equal(t, e.ErrGoodRewardAlreadyExists, err)
	})

	t.Run("error - exists repo fails", func(t *testing.T) {
		mockExistsRepo.
			EXPECT().
			Exists(ctx, testReward.Match).
			Return(false, errors.New("db error"))

		err := rewardService.Register(ctx, testReward)
		assert.Equal(t, e.ErrInternal, err)
	})

	t.Run("error - save repo fails", func(t *testing.T) {
		mockExistsRepo.
			EXPECT().
			Exists(ctx, testReward.Match).
			Return(false, nil)

		mockSaveRepo.
			EXPECT().
			Save(ctx, gomock.Any()).
			Return(errors.New("insert error"))

		err := rewardService.Register(ctx, testReward)
		assert.Equal(t, e.ErrInternal, err)
	})
}
