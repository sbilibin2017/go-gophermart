package middlewares

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type claims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}

func encodeToken(mockConfig *configs.GophermartConfig) (string, error) {
	claims := &claims{
		Login: "test",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(mockConfig.JWTExp)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(mockConfig.JWTSecretKey))
	if err != nil {
		return "", errors.New("err")
	}
	return signedToken, nil
}

func TestAuthMiddleware(t *testing.T) {
	mockConfig := &configs.GophermartConfig{
		JWTSecretKey: "secretkey",
		JWTExp:       time.Hour * 24,
	}
	token, err := encodeToken(mockConfig)
	require.NoError(t, err, "Failed to generate token")
	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	handler := AuthMiddleware(mockConfig)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login := r.Context().Value(loginKey).(string)
		assert.Equal(t, "test", login)
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	mockConfig := &configs.GophermartConfig{
		JWTSecretKey: "secretkey",
		JWTExp:       time.Hour * 24,
	}
	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer invalid_token")
	rr := httptest.NewRecorder()
	handler := AuthMiddleware(mockConfig)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Should not reach here with invalid token")
	}))
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code, "Expected status Unauthorized")
}

func TestAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	mockConfig := &configs.GophermartConfig{
		JWTSecretKey: "secretkey",
		JWTExp:       time.Hour * 24,
	}
	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	handler := AuthMiddleware(mockConfig)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Should not reach here without Authorization header")
	}))
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code, "Expected status Unauthorized")
}
