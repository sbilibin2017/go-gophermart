package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var jwtKey = []byte("test-secret")

func generateValidToken(t *testing.T) string {
	claims := &jwt.RegisteredClaims{
		Subject:   "123",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)
	require.NoError(t, err)
	return tokenStr
}

func newRequestWithAuth(token string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req
}

func TestAuthMiddleware_MissingAuthHeader(t *testing.T) {
	middleware := AuthMiddleware(jwtKey)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Authorization header missing")
}

func TestAuthMiddleware_InvalidHeaderFormat(t *testing.T) {
	middleware := AuthMiddleware(jwtKey)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "InvalidFormat")

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid authorization header format")
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	middleware := AuthMiddleware(jwtKey)

	rr := httptest.NewRecorder()
	req := newRequestWithAuth("invalid.token.value")

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid or expired token")
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	middleware := AuthMiddleware(jwtKey)

	validToken := generateValidToken(t)
	rr := httptest.NewRecorder()
	req := newRequestWithAuth(validToken)

	called := false
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		claims, ok := GetClaims(r)
		require.True(t, ok)
		assert.Equal(t, "123", claims.Subject)
	}))

	handler.ServeHTTP(rr, req)

	assert.True(t, called, "handler should be called")
	assert.Equal(t, http.StatusOK, rr.Code)
}
