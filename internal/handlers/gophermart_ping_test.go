package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGophermartPingHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test case 1: When Pinger is nil, it should return internal server error
	t.Run("when Pinger is nil", func(t *testing.T) {
		handler := GophermartPingHandler(nil)

		// Create a mock request and response recorder
		req, err := http.NewRequest("GET", "/ping", nil)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		// Assert that the status code is 500 (Internal Server Error)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Database connection error")
	})

	// Test case 2: When Ping returns an error, it should return internal server error
	t.Run("when Ping returns an error", func(t *testing.T) {
		mockPinger := NewMockPinger(ctrl)
		mockPinger.EXPECT().Ping().Return(assert.AnError).Times(1)

		handler := GophermartPingHandler(mockPinger)

		// Create a mock request and response recorder
		req, err := http.NewRequest("GET", "/ping", nil)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		// Assert that the status code is 500 (Internal Server Error)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Database connection error")
	})

	// Test case 3: When Ping succeeds, it should return success message
	t.Run("when Ping succeeds", func(t *testing.T) {
		mockPinger := NewMockPinger(ctrl)
		mockPinger.EXPECT().Ping().Return(nil).Times(1)

		handler := GophermartPingHandler(mockPinger)

		// Create a mock request and response recorder
		req, err := http.NewRequest("GET", "/ping", nil)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		// Assert that the status code is 200 (OK)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "Database connection successful")
	})
}
