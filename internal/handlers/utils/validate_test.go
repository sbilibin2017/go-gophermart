package utils

import (
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
