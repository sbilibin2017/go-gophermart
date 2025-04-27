package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type JWTDecoder interface {
	Decode(tokenString string) (map[string]any, error)
}

func AuthMiddleware(
	decoder JWTDecoder,
	jwtSetter func(ctx context.Context, payload map[string]any) context.Context,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr, err := extractAuthHeader(r.Header.Get("Authorization"))
			if err != nil {
				http.Error(w, "Invalid auth header", http.StatusUnauthorized)
				return
			}
			claims, err := decoder.Decode(tokenStr)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}
			ctx := jwtSetter(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
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
