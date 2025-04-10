package middlewares

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/pkg/jwt"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock JWTDecoder
	mockDecoder := NewMockJWTDecoder(ctrl)

	// Create a handler to test
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login, ok := r.Context().Value(loginKey).(string)
		if !ok {
			t.Fatalf("expected login in context, got nil")
		}
		w.Write([]byte("Hello, " + login))
	})

	// Create the middleware with the mock decoder
	middleware := AuthMiddleware(mockDecoder)

	t.Run("Authorization header missing", func(t *testing.T) {
		// Create a request without the Authorization header
		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		assert.NoError(t, err)

		// Record the response
		rr := httptest.NewRecorder()

		// Execute the middleware by passing the handler through it
		handler := middleware(testHandler)
		handler.ServeHTTP(rr, req)

		// Assert the correct error response
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Equal(t, "authorization header missing\n", rr.Body.String())
	})

	t.Run("Invalid Authorization header format", func(t *testing.T) {
		// Create a request with an invalid Authorization header
		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "InvalidToken")

		// Record the response
		rr := httptest.NewRecorder()

		// Execute the middleware by passing the handler through it
		handler := middleware(testHandler)
		handler.ServeHTTP(rr, req)

		// Assert the correct error response
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Equal(t, "invalid authorization header format\n", rr.Body.String())
	})

	t.Run("Invalid JWT token", func(t *testing.T) {
		// Create a request with a valid Authorization header but invalid JWT
		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer invalid_token")

		// Setup the mock decoder to return an error
		mockDecoder.EXPECT().Decode("invalid_token").Return(nil, errors.New("invalid token"))

		// Record the response
		rr := httptest.NewRecorder()

		// Execute the middleware by passing the handler through it
		handler := middleware(testHandler)
		handler.ServeHTTP(rr, req)

		// Assert the correct error response
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Equal(t, "invalid token\n", rr.Body.String())
	})

	t.Run("Valid JWT token", func(t *testing.T) {
		// Create a request with a valid Authorization header
		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer valid_token")

		// Setup the mock decoder to return a valid claim
		mockClaims := &jwt.Claims{Login: "user123"}
		mockDecoder.EXPECT().Decode("valid_token").Return(mockClaims, nil)

		// Record the response
		rr := httptest.NewRecorder()

		// Execute the middleware by passing the handler through it
		handler := middleware(testHandler)
		handler.ServeHTTP(rr, req)

		// Assert the correct response
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "Hello, user123", rr.Body.String())
	})
}
