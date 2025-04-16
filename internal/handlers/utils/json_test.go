package utils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testPayload struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		expected    testPayload
		expectError bool
	}{
		{
			name:        "valid JSON",
			body:        `{"name":"Alice","age":30}`,
			expected:    testPayload{Name: "Alice", Age: 30},
			expectError: false,
		},
		{
			name:        "invalid JSON",
			body:        `{"name":"Bob", "age":}`,
			expectError: true,
		},
		{
			name:        "empty body",
			body:        ``,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.body))
			rec := httptest.NewRecorder()

			var decoded testPayload
			err := Decode(rec, req, &decoded)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, decoded)
			}
		})
	}
}
