package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testRequest struct {
	Name string `json:"name"`
}

type testResponse struct {
	Message string `json:"message"`
}

func TestDecodeRequest(t *testing.T) {
	payload := `{"name":"Sergey"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(payload))
	w := httptest.NewRecorder()
	var data testRequest
	err := decodeRequest(w, req, &data)
	require.NoError(t, err)
	assert.Equal(t, "Sergey", data.Name)
}

func TestDecodeRequest_InvalidJSON(t *testing.T) {
	payload := `{"name":123`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(payload))
	w := httptest.NewRecorder()
	var data testRequest
	err := decodeRequest(w, req, &data)
	require.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}

func TestEncodeResponse(t *testing.T) {
	w := httptest.NewRecorder()
	resp := testResponse{Message: "OK"}
	encodeResponse(w, resp)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	var result testResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "OK", result.Message)
}

func TestWriteTextPlainResponse(t *testing.T) {
	w := httptest.NewRecorder()
	writeTextPlainResponse(w, http.StatusCreated, "created")
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))
	assert.Equal(t, "created", w.Body.String())
}

func TestGetPathParam(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/users/{userID}", func(w http.ResponseWriter, r *http.Request) {
		userID := getPathParam(r, "userID")
		assert.Equal(t, "42", userID)
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest("GET", "/users/42", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestEncodeResponse_Error(t *testing.T) {
	type BadResponse struct {
		BadField chan int `json:"bad"`
	}
	w := httptest.NewRecorder()
	resp := BadResponse{
		BadField: make(chan int),
	}
	encodeResponse(w, resp)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to encode response")
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
}
