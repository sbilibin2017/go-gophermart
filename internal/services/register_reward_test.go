package services

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupRegisterRewardTest(t *testing.T) (
	*gomock.Controller,
	*MockRewardExistsRepository,
	*MockRewardSaveRepository,
	*RegisterRewardSaveService,
) {
	ctrl := gomock.NewController(t)
	mockRewardExistsRepository := NewMockRewardExistsRepository(ctrl)
	mockRewardSaveRepository := NewMockRewardSaveRepository(ctrl)
	service := NewRegisterRewardSaveService(mockRewardExistsRepository, mockRewardSaveRepository)
	return ctrl, mockRewardExistsRepository, mockRewardSaveRepository, service
}

func TestRegisterRewardSaveService_Register(t *testing.T) {
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
			err := service.Register(context.Background(), "testMatch", 100, "testType")
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
