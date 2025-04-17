package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

type RegisterRewardRequest struct {
	Match      string `json:"match" validate:"required"`
	Reward     uint   `json:"reward" validate:"required,gt=0"`
	RewardType string `json:"reward_type" validate:"required,oneof=% pt"`
}

type RegisterRewardResponse struct {
	Message string `json:"message"`
}

func GetRegisterRewardRequestError(err error) error {
	if _, ok := err.(*validator.InvalidValidationError); ok {
		logger.Logger.Error("Ошибка валидации данных:", err)
		return fmt.Errorf("ошибка валидации данных: %v", err)
	}
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			switch e.Tag() {
			case "required":
				logger.Logger.Error(fmt.Sprintf("Поле %s обязательно для заполнения", e.Field()))
				return fmt.Errorf("поле %s обязательно для заполнения", e.Field())
			case "gt":
				logger.Logger.Error(fmt.Sprintf("Поле %s должно быть больше нуля", e.Field()))
				return fmt.Errorf("поле %s должно быть больше нуля", e.Field())
			case "oneof":
				logger.Logger.Error(fmt.Sprintf("Поле %s должно быть одним из следующих значений: %s", e.Field(), e.Param()))
				return fmt.Errorf("поле %s должно быть одним из следующих значений: %s", e.Field(), e.Param())
			}
		}
	}
	logger.Logger.Error("Неизвестная ошибка валидации:", err)
	return nil
}

type RewardExistsFilterDB struct {
	Match string `db:"match"`
}

type RewardDB struct {
	Match      string `db:"match"`
	Reward     uint   `db:"reward"`
	RewardType string `db:"reward_type"`
}
