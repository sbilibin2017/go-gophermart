package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestRewardSaveService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRewardExistsRepo := NewMockRewardExistsRepository(ctrl)
	mockRewardSaveRepo := NewMockRewardSaveRepository(ctrl)

	service := NewRewardSaveService(mockRewardExistsRepo, mockRewardSaveRepo)

	reward := &types.Reward{
		Match:      "test_match",
		Reward:     100,
		RewardType: types.RewardTypePercent,
	}

	tests := []struct {
		name          string
		setupMock     func()
		expectedError error
	}{
		{
			name: "reward already exists",
			setupMock: func() {
				mockRewardExistsRepo.EXPECT().
					Exists(gomock.Any(), gomock.Any()).
					Return(true, nil)
			},
			expectedError: ErrRewardAlreadyExists,
		},
		{
			name: "error checking if reward exists",
			setupMock: func() {
				mockRewardExistsRepo.EXPECT().
					Exists(gomock.Any(), gomock.Any()).
					Return(false, errors.New("some error"))
			},
			expectedError: ErrRewardIsNotRegistered,
		},
		{
			name: "error saving reward",
			setupMock: func() {
				mockRewardExistsRepo.EXPECT().
					Exists(gomock.Any(), gomock.Any()).
					Return(false, nil)
				mockRewardSaveRepo.EXPECT().
					Save(gomock.Any(), gomock.Any()).
					Return(errors.New("some save error"))
			},
			expectedError: ErrRewardIsNotRegistered,
		},
		{
			name: "successful reward registration",
			setupMock: func() {
				mockRewardExistsRepo.EXPECT().
					Exists(gomock.Any(), gomock.Any()).
					Return(false, nil)
				mockRewardSaveRepo.EXPECT().
					Save(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := service.Register(context.Background(), reward)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
