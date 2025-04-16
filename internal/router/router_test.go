package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterHandler_WithMiddlewares(t *testing.T) {
	tests := []struct {
		name         string
		method       HttpMethod
		path         string
		expectedCode int
	}{
		{
			name:         "GET handler with middleware",
			method:       MethodGet,
			path:         "/test",
			expectedCode: http.StatusOK,
		},
		{
			name:         "POST handler with middleware",
			method:       MethodPost,
			path:         "/test",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mainRouter := chi.NewRouter()

			handlerCalled := false
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				w.WriteHeader(http.StatusOK)
			})

			middleware := func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("X-Test-Middleware", "true")
					next.ServeHTTP(w, r)
				})
			}

			RegisterHandler(mainRouter, "/test", tt.method, handler, []func(http.Handler) http.Handler{middleware})

			req, err := http.NewRequest(string(tt.method), "/test", nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			mainRouter.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			assert.True(t, handlerCalled, "handler should be called")
			assert.Equal(t, "true", rr.Header().Get("X-Test-Middleware"), "middleware header should be set")
		})
	}
}
