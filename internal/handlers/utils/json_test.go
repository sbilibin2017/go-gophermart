package utils

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testPayload struct {
	Name string `json:"name"`
}

func TestDecoder_Decode_Success(t *testing.T) {
	body := []byte(`{"name":"test"}`)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))

	var payload testPayload
	decoder := NewDecoder()
	err := decoder.Decode(req, &payload)

	assert.NoError(t, err)
	assert.Equal(t, "test", payload.Name)
}

func TestDecoder_Decode_InvalidJSON(t *testing.T) {
	body := []byte(`{invalid json}`)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))

	var payload testPayload
	decoder := NewDecoder()
	err := decoder.Decode(req, &payload)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrUnprocessableJSON))
}
