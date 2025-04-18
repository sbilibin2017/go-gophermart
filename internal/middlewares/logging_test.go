package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddleware(t *testing.T) {
	handlerOK := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	handlerError := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error response"))
	})

	middlewareOK := LoggingMiddleware(handlerOK)
	reqOK := httptest.NewRequest("GET", "/test", nil)
	rrOK := httptest.NewRecorder()

	middlewareOK.ServeHTTP(rrOK, reqOK)
	assert.Equal(t, http.StatusOK, rrOK.Code)
	assert.Contains(t, rrOK.Body.String(), "test response")

	middlewareError := LoggingMiddleware(handlerError)
	reqError := httptest.NewRequest("GET", "/error", nil)
	rrError := httptest.NewRecorder()

	middlewareError.ServeHTTP(rrError, reqError)
	assert.Equal(t, http.StatusBadRequest, rrError.Code)
	assert.Contains(t, rrError.Body.String(), "error response")

	logger.Info("Testing complete")
}
