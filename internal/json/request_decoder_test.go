package json

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRequestDecoder struct {
	mock.Mock
}

func TestRequestDecoder_Decode(t *testing.T) {
	tests := []struct {
		name           string
		inputJSON      []byte
		expectedResult interface{}
		expectedErr    error
	}{
		{
			name:           "Valid JSON",
			inputJSON:      []byte(`{"name":"John"}`),
			expectedResult: map[string]string{"name": "John"},
			expectedErr:    nil,
		},
		{
			name:        "Invalid JSON",
			inputJSON:   []byte(`{"name": "John"`), // Missing closing brace
			expectedErr: ErrRequestDecoderUnprocessableJSON,
		},
		{
			name:           "Empty JSON",
			inputJSON:      []byte(`{}`),
			expectedResult: map[string]string{},
			expectedErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/", bytes.NewReader(tt.inputJSON))
			w := httptest.NewRecorder()
			decoder := NewRequestDecoder()
			var result map[string]string
			err := decoder.Decode(w, req, &result)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestRequestDecoder_Decode_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest("POST", "/", bytes.NewReader([]byte(`{"name": "John"`)))
	w := httptest.NewRecorder()
	decoder := NewRequestDecoder()
	var result map[string]string
	err := decoder.Decode(w, req, &result)
	assert.ErrorIs(t, err, ErrRequestDecoderUnprocessableJSON)
}

func TestRequestDecoder_Decode_Success(t *testing.T) {
	req := httptest.NewRequest("POST", "/", bytes.NewReader([]byte(`{"name":"Alice"}`)))
	w := httptest.NewRecorder()
	decoder := NewRequestDecoder()
	var result map[string]string
	err := decoder.Decode(w, req, &result)
	assert.NoError(t, err)
	expectedResult := map[string]string{"name": "Alice"}
	assert.Equal(t, expectedResult, result)
}
