package json

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestRequestDecoder_ValidJSON(t *testing.T) {
	body := `{"name":"John", "age":30}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))

	var data TestData
	decoder := NewRequestDecoder()
	err := decoder.Decode(req, &data)

	assert.NoError(t, err)
	assert.Equal(t, "John", data.Name)
	assert.Equal(t, 30, data.Age)
}

func TestRequestDecoder_InvalidJSON(t *testing.T) {
	body := `{"name":"John", "age":}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))

	var data TestData
	decoder := NewRequestDecoder()
	err := decoder.Decode(req, &data)

	assert.Error(t, err)
}

func TestRequestDecoder_EmptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)

	var data TestData
	decoder := NewRequestDecoder()
	err := decoder.Decode(req, &data)

	assert.Error(t, err)
	assert.Equal(t, io.EOF, err)
}
