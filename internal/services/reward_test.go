package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
	e "github.com/sbilibin2017/go-gophermart/internal/errors"
	"github.com/sbilibin2017/go-gophermart/internal/models"
	"github.com/sbilibin2017/go-gophermart/internal/services"
)

func TestRewardService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockExistsRepo := services.NewMockRewardExistsRepository(ctrl)
	mockSaveRepo := services.NewMockRewardSaveRepository(ctrl)

	rewardService := services.NewRewardService(mockExistsRepo, mockSaveRepo)

	testReward := &domain.Reward{
		Match:      "test-match",
		Reward:     100,
		RewardType: domain.Points,
	}

	t.Run("success - new reward", func(t *testing.T) {
		mockExistsRepo.
			EXPECT().
			Exists(ctx, &models.RewardFilter{Match: testReward.Match}).
			Return(false, nil)

		mockSaveRepo.
			EXPECT().
			Save(ctx, gomock.AssignableToTypeOf(&models.RewardDB{})).
			Return(nil)

		err := rewardService.Register(ctx, testReward)
		assert.NoError(t, err)
	})

	t.Run("error - reward already exists", func(t *testing.T) {
		mockExistsRepo.
			EXPECT().
			Exists(ctx, &models.RewardFilter{Match: testReward.Match}).
			Return(true, nil)

		err := rewardService.Register(ctx, testReward)
		assert.Equal(t, e.ErrGoodRewardAlreadyExists, err)
	})

	t.Run("error - exists repo fails", func(t *testing.T) {
		mockExistsRepo.
			EXPECT().
			Exists(ctx, &models.RewardFilter{Match: testReward.Match}).
			Return(false, errors.New("db error"))

		err := rewardService.Register(ctx, testReward)
		assert.Equal(t, e.ErrInternal, err)
	})

	t.Run("error - save repo fails", func(t *testing.T) {
		mockExistsRepo.
			EXPECT().
			Exists(ctx, &models.RewardFilter{Match: testReward.Match}).
			Return(false, nil)

		mockSaveRepo.
			EXPECT().
			Save(ctx, gomock.AssignableToTypeOf(&models.RewardDB{})).
			Return(errors.New("insert error"))

		err := rewardService.Register(ctx, testReward)
		assert.Equal(t, e.ErrInternal, err)
	})
}
