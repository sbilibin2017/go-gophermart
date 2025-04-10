package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SampleRequest struct {
	Name string `json:"name"`
}

type SampleResponse struct {
	Message string `json:"message"`
}

func TestDecodeJSON(t *testing.T) {
	cases := []struct {
		name          string
		body          string
		expectedError bool
		expected      SampleRequest
	}{
		{
			name:          "Valid JSON",
			body:          `{"name": "John"}`,
			expectedError: false,
			expected:      SampleRequest{Name: "John"},
		},
		{
			name:          "Invalid JSON",
			body:          `{"name": "John"`,
			expectedError: true,
		},
		{
			name:          "Empty JSON body",
			body:          ``,
			expectedError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(tc.body))
			w := httptest.NewRecorder()

			var reqBody SampleRequest
			err := DecodeJSON(w, req, &reqBody)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.Name, reqBody.Name)
			}
		})
	}
}

func TestEncodeJSON(t *testing.T) {
	cases := []struct {
		name           string
		statusCode     int
		value          any
		expectError    bool
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Successful response",
			statusCode:     http.StatusOK,
			value:          SampleResponse{Message: "Success"},
			expectError:    false,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Success"}`,
		},
		{
			name:           "Encoding error",
			statusCode:     http.StatusOK,
			value:          make(chan int), // non-serializable type
			expectError:    true,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to write response",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := EncodeJSON(w, tc.statusCode, tc.value)

			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedStatus, w.Result().StatusCode)
				assert.Contains(t, w.Body.String(), tc.expectedBody)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedStatus, w.Result().StatusCode)

				var resBody SampleResponse
				err = json.NewDecoder(w.Body).Decode(&resBody)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedBody, `{"message":"`+resBody.Message+`"}`)
			}
		})
	}
}
