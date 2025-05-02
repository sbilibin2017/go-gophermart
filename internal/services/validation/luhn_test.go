package validation

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestIsLuhnValidationError(t *testing.T) {
	v := validator.New()
	RegisterLuhnValidation(v)
	t.Run("Valid number", func(t *testing.T) {
		type TestStruct struct {
			OrderNumber string `validate:"luhn"`
		}
		validOrderNumber := TestStruct{OrderNumber: "4532015112830366"}
		err := v.Struct(validOrderNumber)
		assert.NoError(t, err)
	})
	t.Run("Invalid number", func(t *testing.T) {
		type TestStruct struct {
			OrderNumber string `validate:"luhn"`
		}
		invalidOrderNumber := TestStruct{OrderNumber: "123456789"}
		err := v.Struct(invalidOrderNumber)
		assert.Error(t, err)
		if err != nil {
			validationErrs := err.(validator.ValidationErrors)
			assert.True(t, IsLuhnValidationError(validationErrs[0]))
		}
	})
	t.Run("Empty number", func(t *testing.T) {
		type TestStruct struct {
			OrderNumber string `validate:"luhn"`
		}
		emptyOrderNumber := TestStruct{OrderNumber: ""}
		err := v.Struct(emptyOrderNumber)
		assert.Error(t, err)
	})
	t.Run("Non-string input", func(t *testing.T) {
		type TestStruct struct {
			OrderNumber int `validate:"luhn"`
		}
		nonStringOrderNumber := TestStruct{OrderNumber: 4532015112830366}
		err := v.Struct(nonStringOrderNumber)
		assert.Error(t, err)
	})
}

func TestLuhnValidation(t *testing.T) {
	t.Run("Valid Luhn number", func(t *testing.T) {
		assert.True(t, validateNumber("4532015112830366"))
	})
	t.Run("Invalid Luhn number", func(t *testing.T) {
		assert.False(t, validateNumber("123456789"))
	})
	t.Run("Empty string", func(t *testing.T) {
		assert.False(t, validateNumber(""))
	})
}
