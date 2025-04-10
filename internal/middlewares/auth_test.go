package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func generateToken(login string, secret string, expireTime time.Duration) (string, error) {
	claims := struct {
		Login string `json:"login"`
		jwt.RegisteredClaims
	}{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "test",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireTime)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func TestAuthMiddleware(t *testing.T) {
	jwtSecret := "mySecret"
	validToken, err := generateToken("testuser", jwtSecret, time.Hour)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	expiredToken, err := generateToken("expireduser", jwtSecret, -time.Hour)
	if err != nil {
		t.Fatalf("Failed to generate expired token: %v", err)
	}

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		expectedLogin  string
	}{
		{
			name:           "Missing authorization header",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid authorization header format",
			token:          "Basic abcdefg",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Valid token",
			token:          validToken,
			expectedStatus: http.StatusOK,
			expectedLogin:  "testuser", // Ожидаем логин из токена
		},
		{
			name:           "Expired token",
			token:          expiredToken,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем middleware с реальным парсером
			middleware := AuthMiddleware(jwtSecret)

			// Создаем тестовый сервер с middleware
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				login, _ := GetLoginFromContext(r.Context())
				assert.Equal(t, tt.expectedLogin, login)
				w.WriteHeader(http.StatusOK)
			})

			// Оборачиваем handler в middleware
			wrappedHandler := middleware(handler)

			// Создаем тестовый запрос
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			// Запускаем тест
			recorder := httptest.NewRecorder()
			wrappedHandler.ServeHTTP(recorder, req)

			// Проверяем статус ответа
			assert.Equal(t, tt.expectedStatus, recorder.Code)
		})
	}
}
