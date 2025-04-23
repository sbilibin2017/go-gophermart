package validators

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/sbilibin2017/go-gophermart/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestRewardTypeValidator(t *testing.T) {
	validate := validator.New()
	RegisterRewardTypeValidator(validate)

	tests := []struct {
		name     string
		input    string
		expected error
	}{
		{
			name:     "Valid reward type - percent",
			input:    string(types.RewardTypePercent),
			expected: nil,
		},
		{
			name:     "Valid reward type - point",
			input:    string(types.RewardTypePoint),
			expected: nil,
		},
		{
			name:     "Invalid reward type",
			input:    "invalid",
			expected: assert.AnError, // Ожидается ошибка для некорректного типа
		},
		{
			name:     "Empty reward type",
			input:    "",
			expected: assert.AnError, // Ожидается ошибка для пустого типа
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Var(tt.input, "reward_type")
			if tt.expected == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "reward_type")
			}
		})
	}
}
