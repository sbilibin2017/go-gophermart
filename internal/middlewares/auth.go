package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type contextKey string

const loginKey contextKey = "login"

var (
	ErrAuthHeaderMissing       = errors.New("authorization header missing")
	ErrInvalidAuthHeaderFormat = errors.New("invalid authorization header format")
)

type JWTConfig interface {
	GetJWTSecretKey() string
}

type Claims interface {
	GetLogin() string
}

func AuthMiddleware(
	config JWTConfig,
	decoder func(tokenStr string, config JWTConfig) (Claims, error),
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := getTokenFromRequestHeader(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}
			claims, err := decoder(tokenString, config)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, loginKey, claims.GetLogin())
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func getTokenFromRequestHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrAuthHeaderMissing
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", ErrInvalidAuthHeaderFormat
	}
	return parts[1], nil
}
