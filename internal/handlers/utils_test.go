package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLoginFromContext_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	mockProvider := func(ctx context.Context) (string, error) {
		return "testuser", nil
	}
	login, err := getLoginFromContext(rr, req, mockProvider)
	require.NoError(t, err)
	assert.Equal(t, "testuser", login)
}

func TestGetLoginFromContext_Error(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	mockProvider := func(ctx context.Context) (string, error) {
		return "", errors.New("unauthorized")
	}
	login, err := getLoginFromContext(rr, req, mockProvider)
	require.Error(t, err)
	assert.Empty(t, login)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Unauthorized")
}

func TestEncodeResponseBody_Success(t *testing.T) {
	rr := httptest.NewRecorder()
	data := map[string]string{"message": "hello"}
	statusCode := http.StatusCreated
	err := encodeResponseBody(rr, data, statusCode)
	require.NoError(t, err)
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	require.Equal(t, statusCode, rr.Code)
	var result map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	require.NoError(t, err)
	require.Equal(t, data, result)
}

func TestEncodeResponseBody_EncodeError(t *testing.T) {
	rr := httptest.NewRecorder()
	err := encodeResponseBody(rr, make(chan int), http.StatusOK)
	require.Error(t, err)
}

func TestGetURLParam(t *testing.T) {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userID", "42")
	req := httptest.NewRequest(http.MethodGet, "/users/42", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	val := getURLParam(req, "userID")
	require.Equal(t, "42", val)
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "Hello"},
		{"HELLO", "Hello"},
		{"h", "H"},
		{"H", "H"},
		{"", ""},
		{"123abc", "123abc"},
		{"aBC", "Abc"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := capitalize(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}
