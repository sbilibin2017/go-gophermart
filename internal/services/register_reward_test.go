package services

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/types"
	"github.com/stretchr/testify/assert"
)

func setupRegisterRewardTest(t *testing.T) (
	*gomock.Controller,
	*MockRegisterRewardExistsRepository,
	*MockRegisterRewardSaveRepository,
	*RegisterRewardService,
) {
	ctrl := gomock.NewController(t)
	mockRewardExistsRepository := NewMockRegisterRewardExistsRepository(ctrl)
	mockRewardSaveRepository := NewMockRegisterRewardSaveRepository(ctrl)
	service := NewRegisterRewardService(mockRewardExistsRepository, mockRewardSaveRepository)
	return ctrl, mockRewardExistsRepository, mockRewardSaveRepository, service
}

func TestRegisterRewardService_Register(t *testing.T) {
	ctrl, mockRewardExistsRepository, mockRewardSaveRepository, service := setupRegisterRewardTest(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		setup         func()
		expectedError error
	}{
		{
			name: "Reward does not exist, save successful",
			setup: func() {
				mockRewardExistsRepository.EXPECT().
					Exists(context.Background(), gomock.Any()).
					Return(false, nil).
					Times(1)

				mockRewardSaveRepository.EXPECT().
					Save(context.Background(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			expectedError: nil,
		},
		{
			name: "Reward already exists",
			setup: func() {
				mockRewardExistsRepository.EXPECT().
					Exists(context.Background(), gomock.Any()).
					Return(true, nil).
					Times(1)
			},
			expectedError: ErrRewardAlreadyExists,
		},
		{
			name: "Error checking if reward exists",
			setup: func() {
				mockRewardExistsRepository.EXPECT().
					Exists(context.Background(), gomock.Any()).
					Return(false, ErrRewardIsNotRegistered).
					Times(1)
			},
			expectedError: ErrRewardIsNotRegistered,
		},
		{
			name: "Error saving reward",
			setup: func() {
				mockRewardExistsRepository.EXPECT().
					Exists(context.Background(), gomock.Any()).
					Return(false, nil).
					Times(1)

				mockRewardSaveRepository.EXPECT().
					Save(context.Background(), gomock.Any()).
					Return(ErrRewardIsNotRegistered).
					Times(1)
			},
			expectedError: ErrRewardIsNotRegistered,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			// Используем правильный тип для RegisterRewardRequest
			req := &types.RegisterRewardRequest{
				Match:      "testMatch",
				Reward:     100,
				RewardType: "testType",
			}

			err := service.Register(context.Background(), req)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
