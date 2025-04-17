package usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/dto"
	"github.com/sbilibin2017/go-gophermart/internal/usecases"
	"github.com/stretchr/testify/assert"
)

func TestRegisterRewardUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := usecases.NewMockRegisterRewardService(ctrl)
	mockValidator := usecases.NewMockRegisterRewardValidator(ctrl)

	uc := usecases.NewRegisterRewardUsecase(mockService, mockValidator)

	req := &dto.RegisterRewardRequest{
		Match:      "product-123",
		Reward:     100,
		RewardType: string(domain.RewardTypePercent),
	}

	tests := []struct {
		name              string
		isErr             bool
		mockExpectActions func()
	}{
		{
			name:  "успешная регистрация",
			isErr: false,
			mockExpectActions: func() {
				mockValidator.EXPECT().
					Struct(req).
					Return(nil).
					Times(1)

				mockService.EXPECT().
					Register(gomock.Any(), &domain.Reward{
						Match:      req.Match,
						Reward:     req.Reward,
						RewardType: domain.RewardType(req.RewardType),
					}).
					Return(nil).
					Times(1)
			},
		},
		{
			name:  "ошибка валидации",
			isErr: true,
			mockExpectActions: func() {
				mockValidator.EXPECT().
					Struct(req).
					Return(errors.New("ошибка валидации")).
					Times(1)
			},
		},
		{
			name:  "ошибка сервиса регистрации",
			isErr: true,
			mockExpectActions: func() {
				mockValidator.EXPECT().
					Struct(req).
					Return(nil).
					Times(1)

				mockService.EXPECT().
					Register(gomock.Any(), &domain.Reward{
						Match:      req.Match,
						Reward:     req.Reward,
						RewardType: domain.RewardType(req.RewardType),
					}).
					Return(errors.New("ошибка сервиса")).
					Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpectActions()

			resp, err := uc.Execute(context.Background(), req)

			if tt.isErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, "Информация о вознаграждении за товар зарегистрирована", resp.Message)
			}
		})
	}
}
