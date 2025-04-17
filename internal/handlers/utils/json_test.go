package utils

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestDecodeJSON_Success(t *testing.T) {
	jsonBody := `{"name":"test","value":42}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonBody))

	var result testStruct
	err := DecodeJSON(req, &result)

	assert.NoError(t, err)
	assert.Equal(t, "test", result.Name)
	assert.Equal(t, 42, result.Value)
}

func TestDecodeJSON_InvalidJSON(t *testing.T) {
	invalidJSON := `this is not json`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(invalidJSON))
	var result testStruct
	err := DecodeJSON(req, &result)
	assert.Error(t, err)

}

func TestDecodeJSON_EmptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	var result testStruct
	err := DecodeJSON(req, &result)
	assert.Error(t, err)
}

func TestEncodeJSON_Success(t *testing.T) {
	recorder := httptest.NewRecorder()
	payload := testStruct{Name: "encoded", Value: 99}
	err := EncodeJSON(recorder, payload)
	assert.NoError(t, err)
	expected := `{"name":"encoded","value":99}` + "\n"
	assert.Equal(t, expected, recorder.Body.String())
}

func TestEncodeJSON_Failure(t *testing.T) {
	brokenWriter := &errorWriter{}
	err := EncodeJSON(brokenWriter, testStruct{Name: "fail", Value: 1})
	assert.Error(t, err)
}

type errorWriter struct{}

func (e *errorWriter) Header() http.Header         { return http.Header{} }
func (e *errorWriter) Write(_ []byte) (int, error) { return 0, errors.New("write error") }
func (e *errorWriter) WriteHeader(statusCode int)  {}
