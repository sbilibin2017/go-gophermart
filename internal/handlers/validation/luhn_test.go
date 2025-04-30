package validation

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestLuhnValidation(t *testing.T) {
	v := validator.New()
	RegisterLuhnValidation(v)

	tests := []struct {
		orderNumber string
		valid       bool
	}{
		{"5062821234567892", true},
		{"5062821234567891", false},
		{"", false},
		{"afasfa", false},
	}

	for _, test := range tests {
		t.Run(test.orderNumber, func(t *testing.T) {
			err := v.Struct(struct {
				OrderNumber string `validate:"luhn"`
			}{
				OrderNumber: test.orderNumber,
			})

			if test.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestRegisterLuhnValidation(t *testing.T) {
	v := validator.New()

	RegisterLuhnValidation(v)

	err := v.Struct(struct {
		OrderNumber string `validate:"luhn"`
	}{
		OrderNumber: "5062821234567892",
	})

	assert.NoError(t, err)
}

func TestIsLuhnValidationError(t *testing.T) {
	v := validator.New()
	RegisterLuhnValidation(v)

	err := v.Struct(struct {
		OrderNumber int `validate:"luhn"`
	}{
		OrderNumber: 79927398714, // Invalid Luhn number
	})

	if assert.Error(t, err) {
		fieldErr := err.(validator.ValidationErrors)[0]
		assert.True(t, IsLuhnValidationError(fieldErr))
	}
}
