package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const claimsContextKey contextKey = "claims"

var (
	ErrMissingAuthHeader     = errors.New("authorization header missing")
	ErrInvalidAuthHeader     = errors.New("invalid authorization header format")
	ErrInvalidOrExpiredToken = errors.New("invalid or expired token")
)

func AuthMiddleware(jwtKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, ErrMissingAuthHeader.Error(), http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, ErrInvalidAuthHeader.Error(), http.StatusUnauthorized)
				return
			}

			claims, err := parseToken(parts[1], jwtKey)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), claimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func parseToken(tokenStr string, jwtKey []byte) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidOrExpiredToken
	}

	return claims, nil
}

func GetClaims(r *http.Request) (*jwt.RegisteredClaims, bool) {
	claims, ok := r.Context().Value(claimsContextKey).(*jwt.RegisteredClaims)
	return claims, ok
}
