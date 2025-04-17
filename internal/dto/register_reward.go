package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
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
		return fmt.Errorf("ошибка валидации данных: %v", err)
	}
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			switch e.Tag() {
			case "required":
				return fmt.Errorf("поле %s обязательно для заполнения", e.Field())
			case "gt":
				return fmt.Errorf("поле %s должно быть больше нуля", e.Field())
			case "oneof":
				return fmt.Errorf("поле %s должно быть одним из следующих значений: %s", e.Field(), e.Param())
			}
		}
	}
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
