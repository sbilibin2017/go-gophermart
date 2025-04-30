package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddleware_BasicResponse(t *testing.T) {
	logger.InitLoggerWithInfoLevel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("hello"))
	})

	loggedHandler := middlewares.LoggingMiddleware(handler)

	req := httptest.NewRequest("GET", "/test-uri", nil)
	rec := httptest.NewRecorder()

	start := time.Now()
	loggedHandler.ServeHTTP(rec, req)
	duration := time.Since(start)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "hello", rec.Body.String())

	assert.LessOrEqual(t, duration.Milliseconds(), int64(100), "Duration should be small in test")
}
