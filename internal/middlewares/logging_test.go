package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestLoggingMiddleware(t *testing.T) {

	logger.Init(zapcore.InfoLevel)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World"))
	})

	middleware := LoggingMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "Hello, World", rr.Body.String())

}

func TestResponseWriter(t *testing.T) {
	// Создаем mock http.ResponseWriter
	rr := httptest.NewRecorder()
	rw := &ResponseWriter{ResponseWriter: rr}

	// Тестируем WriteHeader
	rw.WriteHeader(http.StatusNotFound)
	assert.Equal(t, http.StatusNotFound, rw.statusCode)

	// Тестируем Write
	_, err := rw.Write([]byte("Test body"))
	require.NoError(t, err)
	assert.Equal(t, 9, rw.responseSize)
}
