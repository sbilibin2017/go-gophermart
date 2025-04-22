package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

const (
	ErrMissingAuthHeader     = "Authorization header missing"
	ErrInvalidAuthHeader     = "Invalid authorization header format"
	ErrInvalidOrExpiredToken = "Invalid or expired token"
)

func AuthMiddleware(jwtKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, ErrMissingAuthHeader, http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, ErrInvalidAuthHeader, http.StatusUnauthorized)
				return
			}

			claims, err := parseToken(parts[1], jwtKey)
			if err != nil {
				http.Error(w, ErrInvalidOrExpiredToken, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), claimsContextKey, claims.Login)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetLogin(r *http.Request) (string, bool) {
	login, ok := r.Context().Value(claimsContextKey).(string)
	return login, ok
}

type claims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}

type contextJWTKey string

const claimsContextKey contextJWTKey = "login"

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
