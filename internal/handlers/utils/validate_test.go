package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/sbilibin2017/go-gophermart/internal/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Пример структуры для валидации
type TestStruct struct {
	Name string `validate:"required"`
	Age  int    `validate:"gte=0,lte=130"`
}

func TestValidateStruct(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name           string
		input          interface{}
		expectedStatus int
		expectError    bool
	}{
		{
			name: "valid struct",
			input: TestStruct{
				Name: "Alice",
				Age:  30,
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "missing name",
			input: TestStruct{
				Name: "",
				Age:  25,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "invalid age",
			input: TestStruct{
				Name: "Bob",
				Age:  -5,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			err := ValidateStruct(rec, validate, tt.input)

			if tt.expectError {
				require.Error(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
				assert.Equal(t, errors.ErrDataIsNotValid.Error()+"\n", rec.Body.String())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
				assert.Empty(t, rec.Body.String())
			}
		})
	}
}
