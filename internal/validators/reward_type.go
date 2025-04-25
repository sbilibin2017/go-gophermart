package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/sbilibin2017/go-gophermart/internal/constants"
)

func RegisterRewardTypeValidator(validate *validator.Validate) {
	validate.RegisterValidation("reward_type", rewardTypeValidator)
}

func rewardTypeValidator(fl validator.FieldLevel) bool {
	rewardType := fl.Field().String()
	return validateRewardType(rewardType)
}

func validateRewardType(rewardType string) bool {
	if rewardType == constants.REWARD_TYPE_PERCENT || rewardType == constants.REWARD_TYPE_POINT {
		return true
	}
	return false
}
