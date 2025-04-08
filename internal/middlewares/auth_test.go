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
)

func TestGetTokenHeader(t *testing.T) {
	tests := []struct {
		name          string
		authHeader    string
		expectedToken string
		expectedError error
	}{
		{
			name:          "Valid token header",
			authHeader:    "Bearer validtoken123",
			expectedToken: "validtoken123",
			expectedError: nil,
		},
		{
			name:          "Missing token",
			authHeader:    "",
			expectedToken: "",
			expectedError: errors.New("authorization header missing"),
		},
		{
			name:          "Invalid token format",
			authHeader:    "InvalidFormat token",
			expectedToken: "",
			expectedError: errors.New("invalid authorization format"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", tt.authHeader)
			token, err := getTokenHeader(req)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedToken, token)
			}
		})
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	config := &configs.JWTConfig{
		SecretKey: "secret",
		Exp:       time.Hour,
	}
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login := r.Context().Value(userLoginKey).(string)
		w.Write([]byte("Authenticated user: " + login))
	})
	authMiddleware := AuthMiddleware(config)(testHandler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken123")
	rr := httptest.NewRecorder()
	authMiddleware.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "Invalid or expired token\n", rr.Body.String())
}

func TestAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	config := &configs.JWTConfig{
		SecretKey: "secret",
		Exp:       time.Hour,
	}
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called when token is missing")
	})
	authMiddleware := AuthMiddleware(config)(testHandler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	authMiddleware.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "Invalid token header\n", rr.Body.String())
}

func TestAuthMiddleware_TokenValidationFailure(t *testing.T) {
	config := &configs.JWTConfig{
		SecretKey: "secret",
		Exp:       time.Hour,
	}
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called when token is invalid")
	})
	authMiddleware := AuthMiddleware(config)(testHandler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken123")
	rr := httptest.NewRecorder()
	authMiddleware.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "Invalid or expired token\n", rr.Body.String())
}

func TestAuthMiddleware_ContextWithLogin(t *testing.T) {
	config := &configs.JWTConfig{
		SecretKey: "secret",
		Exp:       3600,
	}
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login, ok := r.Context().Value(userLoginKey).(string)
		if !ok {
			t.Fatal("Failed to retrieve login from context")
		}
		w.Write([]byte("Authenticated user: " + login))
	})
	authMiddleware := AuthMiddleware(config)(testHandler)
	claims := jwt.MapClaims{
		"login": "testuser",
		"exp":   time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		t.Fatal("Error signing token:", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rr := httptest.NewRecorder()
	authMiddleware.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Authenticated user: testuser", rr.Body.String())
}
