package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/sbilibin2017/go-gophermart/internal/jwt"
)

type contextKey string

const (
	loginKey contextKey = "login"
)

var (
	ErrAuthHeaderMissing       = errors.New("authorization header missing")
	ErrInvalidAuthHeaderFormat = errors.New("invalid authorization header format")
)

type JWTPArser interface {
	Parse(tokenStr string) (*jwt.Claims, error)
}

func AuthMiddleware(p JWTPArser) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := getTokenFromRequestHeader(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			claims, err := p.Parse(tokenString)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), loginKey, claims)
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

func GetLoginFromContext(ctx context.Context) (string, bool) {
	claims, ok := ctx.Value(loginKey).(*jwt.Claims)
	if !ok || claims == nil {
		return "", false
	}
	return claims.Login, true
}
