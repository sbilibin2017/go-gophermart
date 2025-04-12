package json

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseEncoder_Encode_Valid(t *testing.T) {
	rr := httptest.NewRecorder()
	encoder := NewResponseEncoder()

	data := map[string]string{"message": "ok"}
	err := encoder.Encode(rr, http.StatusOK, data)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"message":"ok"}`, rr.Body.String())
}

func TestResponseEncoder_Encode_NilBody(t *testing.T) {
	rr := httptest.NewRecorder()
	encoder := NewResponseEncoder()

	err := encoder.Encode(rr, http.StatusNoContent, nil)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Empty(t, rr.Body.String())
}

func TestResponseEncoder_Encode_InvalidType(t *testing.T) {
	rr := httptest.NewRecorder()
	encoder := NewResponseEncoder()

	// Канал нельзя сериализовать в JSON
	invalid := make(chan int)
	err := encoder.Encode(rr, http.StatusOK, invalid)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "json: unsupported type")
}
