package validators

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func init() {
	logger.Init(zapcore.InfoLevel)
}

func TestLuhnCheck(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("luhn", ValidateLuhn)

	tests := []struct {
		number   string
		expected bool
	}{
		{"4539 1488 0343 6467", true},  // Корректный номер
		{"6011514433546201", true},     // Корректный номер
		{"123456789", false},           // Некорректный номер
		{"0000 0000 0000 0000", true},  // Корректный номер
		{"1234 5678 9012 3456", false}, // Некорректный номер
		{"abcd efgh ijkl mnop", false}, // Некорректный номер (не цифры)
	}

	for _, test := range tests {
		t.Run(test.number, func(t *testing.T) {
			err := validate.Var(test.number, "luhn")
			if test.expected {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
