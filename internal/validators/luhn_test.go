package validators

import (
	"testing"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func init() {
	logger.Init(zapcore.InfoLevel)
}

func TestLuhnCheck(t *testing.T) {
	tests := []struct {
		number   string
		expected bool
	}{
		{"4539 1488 0343 6467", true},
		{"6011514433546201", true},
		{"123456789", false},
		{"0000 0000 0000 0000", true},
		{"1234 5678 9012 3456", false},
		{"abcd efgh ijkl mnop", false},
	}

	for _, test := range tests {
		t.Run(test.number, func(t *testing.T) {
			result := ValidateNumberWithLuhn(test.number)
			assert.Equal(t, test.expected, result)
		})
	}
}
