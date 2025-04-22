package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

func RegisterRewardTypeValidator(validate *validator.Validate) {
	validate.RegisterValidation("reward_type", rewardTypeValidator)
}

func rewardTypeValidator(fl validator.FieldLevel) bool {
	rewardType := fl.Field().String()
	return validateRewardType(rewardType)
}

func validateRewardType(rewardType string) bool {
	if rewardType == string(types.RewardTypePercent) || rewardType == string(types.RewardTypePoint) {
		return true
	}
	return false
}
