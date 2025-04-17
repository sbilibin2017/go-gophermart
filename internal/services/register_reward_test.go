package services_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/dto"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterRewardService_Register(t *testing.T) {
	tests := []struct {
		name              string
		isErr             bool
		mockExpectActions func(mockRewardExistsRepo *services.MockRewardExistsRepository, mockRewardSaveRepo *services.MockRewardSaveRepository, mockTx *services.MockTx)
		expectedError     error
	}{
		{
			name:  "Награда существует, ошибка",
			isErr: true,
			mockExpectActions: func(mockRewardExistsRepo *services.MockRewardExistsRepository, mockRewardSaveRepo *services.MockRewardSaveRepository, mockTx *services.MockTx) {
				mockRewardExistsRepo.EXPECT().
					Exists(gomock.Any(), gomock.Eq(&dto.RewardExistsFilterDB{Match: "some_match"})).
					Return(true, nil).Times(1)
				mockRewardSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)
				mockTx.EXPECT().
					Do(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
						return fn(nil)
					}).Times(1)
			},
			expectedError: services.ErrGoodRewardAlreadyExists,
		},
		{
			name:  "Ошибка при проверке существования награды",
			isErr: true,
			mockExpectActions: func(mockRewardExistsRepo *services.MockRewardExistsRepository, mockRewardSaveRepo *services.MockRewardSaveRepository, mockTx *services.MockTx) {
				mockRewardExistsRepo.EXPECT().
					Exists(gomock.Any(), gomock.Eq(&dto.RewardExistsFilterDB{Match: "some_match"})).
					Return(false, assert.AnError).Times(1)
				mockRewardSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)
				mockTx.EXPECT().
					Do(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
						return fn(nil)
					}).Times(1)
			},
			expectedError: assert.AnError,
		},
		{
			name:  "Награда не существует, сохранение успешно",
			isErr: false,
			mockExpectActions: func(mockRewardExistsRepo *services.MockRewardExistsRepository, mockRewardSaveRepo *services.MockRewardSaveRepository, mockTx *services.MockTx) {
				mockRewardExistsRepo.EXPECT().
					Exists(gomock.Any(), gomock.Eq(&dto.RewardExistsFilterDB{Match: "some_match"})).
					Return(false, nil).Times(1)
				mockRewardSaveRepo.EXPECT().
					Save(gomock.Any(), gomock.Eq(&dto.RewardDB{
						Match:      "some_match",
						Reward:     100,
						RewardType: "bonus",
					})).
					Return(nil).Times(1)
				mockTx.EXPECT().
					Do(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
						return fn(nil)
					}).Times(1)
			},
			expectedError: nil,
		},
		{
			name:  "Награда существует, сохранение не должно быть вызвано",
			isErr: true,
			mockExpectActions: func(mockRewardExistsRepo *services.MockRewardExistsRepository, mockRewardSaveRepo *services.MockRewardSaveRepository, mockTx *services.MockTx) {
				mockRewardExistsRepo.EXPECT().
					Exists(gomock.Any(), gomock.Eq(&dto.RewardExistsFilterDB{Match: "some_match"})).
					Return(true, nil).Times(1)
				mockRewardSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)
				mockTx.EXPECT().
					Do(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
						return fn(nil)
					}).Times(1)
			},
			expectedError: services.ErrGoodRewardAlreadyExists,
		},
		{
			name:  "Ошибка при сохранении, возвращение ошибки",
			isErr: true,
			mockExpectActions: func(mockRewardExistsRepo *services.MockRewardExistsRepository, mockRewardSaveRepo *services.MockRewardSaveRepository, mockTx *services.MockTx) {
				mockRewardExistsRepo.EXPECT().
					Exists(gomock.Any(), gomock.Eq(&dto.RewardExistsFilterDB{Match: "some_match"})).
					Return(false, nil).Times(1)
				mockRewardSaveRepo.EXPECT().
					Save(gomock.Any(), gomock.Eq(&dto.RewardDB{
						Match:      "some_match",
						Reward:     100,
						RewardType: "bonus",
					})).
					Return(assert.AnError).Times(1)
				mockTx.EXPECT().
					Do(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
						return fn(nil)
					}).Times(1)
			},
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRewardExistsRepo := services.NewMockRewardExistsRepository(ctrl)
			mockRewardSaveRepo := services.NewMockRewardSaveRepository(ctrl)
			mockTx := services.NewMockTx(ctrl)

			tt.mockExpectActions(mockRewardExistsRepo, mockRewardSaveRepo, mockTx)

			service := services.NewRegisterRewardService(mockRewardExistsRepo, mockRewardSaveRepo, mockTx)

			reward := &domain.Reward{
				Match:      "some_match",
				Reward:     100,
				RewardType: domain.RewardType("bonus"),
			}

			err := service.Register(context.Background(), reward)

			if tt.isErr {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
