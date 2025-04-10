package utils

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestHandleValidationError(t *testing.T) {
	w := httptest.NewRecorder()
	t.Run("Internal Server Error for Non-Validation Error", func(t *testing.T) {
		HandleValidationError(w, errors.New("some internal error"))
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		assert.Contains(t, w.Body.String(), "Internal Server Error")
	})
}

func TestHandleValidationError_FirstErrorOnly(t *testing.T) {
	validate := validator.New()
	type SampleRequest struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}
	request := SampleRequest{
		Name:  "",
		Email: "invalid-email",
	}
	err := validate.Struct(request)
	assert.Error(t, err)
	validationErrors, ok := err.(validator.ValidationErrors)
	assert.True(t, ok)
	w := httptest.NewRecorder()
	HandleValidationError(w, validationErrors)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "Name")
	assert.NotContains(t, w.Body.String(), "Email")
}
