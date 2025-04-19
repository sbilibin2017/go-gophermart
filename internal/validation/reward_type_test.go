package validation

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/sbilibin2017/go-gophermart/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestValidateRewardType(t *testing.T) {
	validate := validator.New()

	validate.RegisterValidation("reward_type", ValidateRewardType)

	tests := []struct {
		name      string
		input     interface{}
		expectErr bool
	}{
		{
			name:      "Valid RewardTypePercent",
			input:     types.RewardTypePercent,
			expectErr: false,
		},
		{
			name:      "Valid RewardTypePoints",
			input:     types.RewardTypePoints,
			expectErr: false,
		},
		{
			name:      "Invalid RewardType (not a valid type)",
			input:     "invalid_type",
			expectErr: true,
		},
		{
			name:      "Invalid RewardType (empty value)",
			input:     nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Var(tt.input, "reward_type")
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
