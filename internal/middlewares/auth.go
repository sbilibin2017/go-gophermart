package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := getTokenString(r)
			if err != nil {
				handleUnauthorizedError(w)
				return
			}
			login, err := getLoginFromToken(tokenString, secretKey)
			if err != nil {
				handleUnauthorizedError(w)
				return
			}
			ctx := setLoginToContext(r.Context(), login)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func handleUnauthorizedError(w http.ResponseWriter) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func getTokenString(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("authorization header is missing or invalid")
	}
	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

type contextLoginKey string

const loginContextKey contextLoginKey = "login"

func setLoginToContext(ctx context.Context, login string) context.Context {
	return context.WithValue(ctx, loginContextKey, login)
}

func GetLoginFromContext(ctx context.Context) (string, bool) {
	login, ok := ctx.Value(loginContextKey).(string)
	return login, ok
}

type claims struct {
	jwt.RegisteredClaims
	Login string `json:"login"`
}

func getLoginFromToken(tokenString string, secretKey string) (string, error) {
	claims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (any, error) {
			return []byte(secretKey), nil
		})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}
	if !token.Valid {
		return "", fmt.Errorf("token is not valid")
	}
	return claims.Login, nil
}

func GenerateTokenString(login string, secretKey string, exp time.Duration) (string, error) {
	claims := claims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "gophermart",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}
	return signedToken, nil
}
