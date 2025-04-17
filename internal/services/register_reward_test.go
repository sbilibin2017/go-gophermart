package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (*gomock.Controller, *MockRewardExistsRepository, *MockRewardSaveRepository, *RegisterRewardService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockExists := NewMockRewardExistsRepository(ctrl)
	mockSave := NewMockRewardSaveRepository(ctrl)
	service := NewRegisterRewardService(mockExists, mockSave, nil) // Теперь не нужен DB
	return ctrl, mockExists, mockSave, service
}

func TestRegister(t *testing.T) {
	type mockBehavior func(
		mockExists *MockRewardExistsRepository,
		mockSave *MockRewardSaveRepository,
	)

	tests := []struct {
		name          string
		reward        *domain.Reward
		mockBehavior  mockBehavior
		expectedError error
	}{
		{
			name: "Success",
			reward: &domain.Reward{
				Match:      "order123",
				Reward:     100,
				RewardType: domain.RewardTypePoints,
			},
			mockBehavior: func(mockExists *MockRewardExistsRepository, mockSave *MockRewardSaveRepository) {
				mockExists.EXPECT().Exists(gomock.Any(), map[string]any{"match": "order123"}).Return(false, nil)
				mockSave.EXPECT().Save(gomock.Any(), map[string]any{
					"match":       "order123",
					"reward":      uint(100),
					"reward_type": "pt",
				}).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "AlreadyExists",
			reward: &domain.Reward{
				Match:      "order123",
				Reward:     100,
				RewardType: domain.RewardTypePoints,
			},
			mockBehavior: func(mockExists *MockRewardExistsRepository, _ *MockRewardSaveRepository) {
				mockExists.EXPECT().Exists(gomock.Any(), map[string]any{"match": "order123"}).Return(true, nil)
			},
			expectedError: domain.ErrRewardKeyAlreadyRegistered,
		},
		{
			name: "ExistsReturnsError",
			reward: &domain.Reward{
				Match:      "order123",
				Reward:     100,
				RewardType: domain.RewardTypePoints,
			},
			mockBehavior: func(mockExists *MockRewardExistsRepository, _ *MockRewardSaveRepository) {
				mockExists.EXPECT().Exists(gomock.Any(), map[string]any{"match": "order123"}).Return(false, errors.New("db error"))
			},
			expectedError: domain.ErrFailedToRegisterReward,
		},
		{
			name: "SaveReturnsError",
			reward: &domain.Reward{
				Match:      "order123",
				Reward:     100,
				RewardType: domain.RewardTypePoints,
			},
			mockBehavior: func(mockExists *MockRewardExistsRepository, mockSave *MockRewardSaveRepository) {
				mockExists.EXPECT().Exists(gomock.Any(), map[string]any{"match": "order123"}).Return(false, nil)
				mockSave.EXPECT().Save(gomock.Any(), map[string]any{
					"match":       "order123",
					"reward":      uint(100),
					"reward_type": "pt",
				}).Return(errors.New("insert error"))
			},
			expectedError: domain.ErrFailedToRegisterReward,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, mockExists, mockSave, service := setup(t)
			defer ctrl.Finish()
			tt.mockBehavior(mockExists, mockSave)
			err := service.Register(context.Background(), tt.reward)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
