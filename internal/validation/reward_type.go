package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

func ValidateRewardType(fl validator.FieldLevel) bool {
	rt, ok := fl.Field().Interface().(types.RewardType)
	if !ok {
		return false
	}
	return rt == types.RewardTypePercent || rt == types.RewardTypePoints
}
