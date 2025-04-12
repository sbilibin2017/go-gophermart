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

	t.Run("when Pinger is nil", func(t *testing.T) {
		handler := PingHandler(nil)

		req, err := http.NewRequest("GET", "/ping", nil)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Database connection error")
	})

	t.Run("when Ping returns an error", func(t *testing.T) {
		mockPinger := NewMockPinger(ctrl)
		mockPinger.EXPECT().Ping().Return(assert.AnError).Times(1)

		handler := PingHandler(mockPinger)

		req, err := http.NewRequest("GET", "/ping", nil)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Database connection error")
	})

	t.Run("when Ping succeeds", func(t *testing.T) {
		mockPinger := NewMockPinger(ctrl)
		mockPinger.EXPECT().Ping().Return(nil).Times(1)

		handler := PingHandler(mockPinger)

		req, err := http.NewRequest("GET", "/ping", nil)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "Database connection successful")
	})
}
