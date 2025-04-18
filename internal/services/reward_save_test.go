package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

func TestRegisterRewardService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	reward := &types.Reward{
		Match:      "test-match",
		Reward:     100,
		RewardType: types.RewardType("bonus"),
	}

	tests := []struct {
		name        string
		mockSetup   func(re *services.MockRewardExistsRepository, rs *services.MockRewardSaveRepository)
		expectedErr error
	}{
		{
			name: "success - reward does not exist and saved",
			mockSetup: func(re *services.MockRewardExistsRepository, rs *services.MockRewardSaveRepository) {
				re.EXPECT().
					Exists(ctx, map[string]any{"match": reward.Match}).
					Return(false, nil)
				rs.EXPECT().
					Save(ctx, map[string]any{
						"match":       reward.Match,
						"reward":      reward.Reward,
						"reward_type": string(reward.RewardType),
					}).
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "error - reward already exists",
			mockSetup: func(re *services.MockRewardExistsRepository, rs *services.MockRewardSaveRepository) {
				re.EXPECT().
					Exists(ctx, map[string]any{"match": reward.Match}).
					Return(true, nil)
			},
			expectedErr: services.ErrRewardAlreadyExists,
		},
		{
			name: "error - exists check fails",
			mockSetup: func(re *services.MockRewardExistsRepository, rs *services.MockRewardSaveRepository) {
				re.EXPECT().
					Exists(ctx, map[string]any{"match": reward.Match}).
					Return(false, errors.New("db error"))
			},
			expectedErr: services.ErrRewardIsNotRegistered,
		},
		{
			name: "error - save fails",
			mockSetup: func(re *services.MockRewardExistsRepository, rs *services.MockRewardSaveRepository) {
				re.EXPECT().
					Exists(ctx, map[string]any{"match": reward.Match}).
					Return(false, nil)
				rs.EXPECT().
					Save(ctx, map[string]any{
						"match":       reward.Match,
						"reward":      reward.Reward,
						"reward_type": string(reward.RewardType),
					}).
					Return(errors.New("save error"))
			},
			expectedErr: services.ErrRewardIsNotRegistered,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExistsRepo := services.NewMockRewardExistsRepository(ctrl)
			mockSaveRepo := services.NewMockRewardSaveRepository(ctrl)

			tt.mockSetup(mockExistsRepo, mockSaveRepo)

			service := services.NewRewardSaveService(mockExistsRepo, mockSaveRepo)
			err := service.Register(ctx, reward)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
