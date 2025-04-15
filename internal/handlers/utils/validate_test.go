package utils

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
}

func TestValidate_Success(t *testing.T) {
	validObj := testStruct{
		Name:  "John Doe",
		Email: "johndoe@example.com",
	}
	v := validator.New()
	err := Validate(v, validObj)
	assert.NoError(t, err)
}

func TestValidate_Failure(t *testing.T) {
	invalidObj := testStruct{
		Name:  "",
		Email: "not-an-email",
	}
	v := validator.New()
	err := Validate(v, invalidObj)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
	assert.Contains(t, err.Error(), "Name")
	assert.Contains(t, err.Error(), "Email")
}

type mockValidatorNonValidationError struct{}

func (m mockValidatorNonValidationError) Struct(i interface{}) error {
	return errors.New("some unexpected error")
}

type mockValidatorValidationError struct{}

func (m mockValidatorValidationError) Struct(i interface{}) error {
	v := validator.New()
	type testStruct struct {
		Field string `validate:"required"`
	}
	return v.Struct(testStruct{})
}

func TestValidate_ReturnsWrappedValidationError(t *testing.T) {
	v := mockValidatorValidationError{}
	err := Validate(v, struct{}{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed:")
}

func TestValidate_ReturnsNonValidationError(t *testing.T) {
	v := mockValidatorNonValidationError{}
	err := Validate(v, struct{}{})
	assert.EqualError(t, err, "some unexpected error")
}

func TestValidate_ReturnsNil(t *testing.T) {
	v := validator.New()

	type validStruct struct {
		Name string `validate:"required"`
	}

	err := Validate(v, validStruct{Name: "test"})
	assert.NoError(t, err)
}
