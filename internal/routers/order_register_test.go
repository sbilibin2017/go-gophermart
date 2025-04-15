package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func init() {
	logger.Init(zapcore.InfoLevel)
}

func TestRegisterOrderRegisterRoute(t *testing.T) {
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	mockGzipMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
	mockLoggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}

	mockRouter := chi.NewRouter()

	RegisterOrderRegisterRoute(mockRouter, "/api", mockHandler, mockGzipMiddleware, mockLoggingMiddleware)

	req, err := http.NewRequest("POST", "/api/orders", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	mockRouter.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, http.StatusOK, rr.Code)
}
