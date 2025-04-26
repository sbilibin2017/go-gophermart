package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(jwtKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Извлекаем и проверяем токен из заголовка Authorization
			tokenStr, err := extractAuthHeader(r.Header.Get("Authorization"))
			if err != nil {
				http.Error(w, "Invalid auth header", http.StatusUnauthorized)
				return
			}

			claims, err := parseToken(tokenStr, jwtKey)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), loginContextKey, claims.Login)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GenerateToken - генерирует JWT токен для пользователя
func GenerateToken(login string, jwtKey []byte, expiresAt time.Duration) (string, error) {
	// Создаем токен
	claims := &claims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
		},
	}

	// Создаем токен с использованием алгоритма HMAC
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GetLogin(ctx context.Context) *string {
	login, ok := ctx.Value(loginContextKey).(string)
	if !ok {
		return nil
	}
	return &login
}

type claims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}

func parseToken(tokenStr string, jwtKey []byte) (*claims, error) {
	claims := &claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}

func extractAuthHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header format")
	}
	return parts[1], nil
}

type contextLoginKey string

const loginContextKey contextLoginKey = "login"
