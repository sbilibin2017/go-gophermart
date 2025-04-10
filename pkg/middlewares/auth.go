package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/sbilibin2017/go-gophermart/pkg/jwt"
)

type contextKey string

const loginKey contextKey = "login"

var (
	ErrAuthHeaderMissing       = errors.New("authorization header missing")
	ErrInvalidAuthHeaderFormat = errors.New("invalid authorization header format")
)

type JWTDecoder interface {
	Decode(tokenStr string) (*jwt.Claims, error)
}

type AuthMiddleware struct {
	Decoder JWTDecoder
}

func (am *AuthMiddleware) Apply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := getTokenFromRequestHeader(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		claims, err := am.Decoder.Decode(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), loginKey, claims.Login)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
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
