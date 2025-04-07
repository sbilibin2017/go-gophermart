package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/jwt"
)

type contextKey string

const userLoginKey contextKey = "userLogin"

func AuthMiddleware(config *configs.JWTConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header missing", http.StatusUnauthorized)
				return
			}
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
				return
			}
			tokenString := tokenParts[1]
			login, err := jwt.GetUserLogin(config, tokenString)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, userLoginKey, login)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
