package utils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}

func TestDecodeJSONBody(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		expectedError error
		expectedData  *testStruct
	}{
		{
			name:          "valid JSON",
			body:          `{"field1": "value", "field2": 123}`,
			expectedError: nil,
			expectedData:  &testStruct{Field1: "value", Field2: 123},
		},
		{
			name:          "invalid JSON",
			body:          `{"field1": "value", "field2":}`,
			expectedError: ErrInvalidRequestBody,
			expectedData:  nil,
		},
		{
			name:          "empty body",
			body:          ``,
			expectedError: ErrInvalidRequestBody,
			expectedData:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()

			var result testStruct
			err := DecodeJSONBody(w, req, &result)

			assert.Equal(t, tt.expectedError, err)
			if tt.expectedData != nil {
				assert.Equal(t, *tt.expectedData, result)
			}
		})
	}
}
