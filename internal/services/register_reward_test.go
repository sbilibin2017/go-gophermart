package services

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
)

func setup(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, *gomock.Controller, *MockRewardExistsRepository, *MockRewardSaveRepository, *RegisterRewardService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockExists := NewMockRewardExistsRepository(ctrl)
	mockSave := NewMockRewardSaveRepository(ctrl)
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	service := NewRegisterRewardService(mockExists, mockSave, sqlxDB)
	return sqlxDB, mockDB, ctrl, mockExists, mockSave, service
}

func TestRegister(t *testing.T) {
	type mockBehavior func(
		mockExists *MockRewardExistsRepository,
		mockSave *MockRewardSaveRepository,
		mockDB sqlmock.Sqlmock,
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
			mockBehavior: func(mockExists *MockRewardExistsRepository, mockSave *MockRewardSaveRepository, mockDB sqlmock.Sqlmock) {
				mockDB.ExpectBegin()
				mockDB.ExpectCommit()

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
			mockBehavior: func(mockExists *MockRewardExistsRepository, _ *MockRewardSaveRepository, mockDB sqlmock.Sqlmock) {
				mockDB.ExpectBegin()
				mockDB.ExpectRollback()

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
			mockBehavior: func(mockExists *MockRewardExistsRepository, _ *MockRewardSaveRepository, mockDB sqlmock.Sqlmock) {
				mockDB.ExpectBegin()
				mockDB.ExpectRollback()

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
			mockBehavior: func(mockExists *MockRewardExistsRepository, mockSave *MockRewardSaveRepository, mockDB sqlmock.Sqlmock) {
				mockDB.ExpectBegin()
				mockDB.ExpectRollback()

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
			sqlxDB, mockDB, ctrl, mockExists, mockSave, service := setup(t)
			defer ctrl.Finish()
			tt.mockBehavior(mockExists, mockSave, mockDB)
			err := service.Register(context.Background(), tt.reward)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mockDB.ExpectationsWereMet())
			_ = sqlxDB.Close()
		})
	}
}
