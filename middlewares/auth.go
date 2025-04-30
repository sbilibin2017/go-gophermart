package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/jwt"
)

func AuthMiddleware(config *configs.JWTConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := getTokenStringFromHeader(r)
			if err != nil {
				handleUnauthorizedError(w)
				return
			}
			claims, err := jwt.GetClaims(tokenString, config.SecretKey)
			if err != nil {
				handleUnauthorizedError(w)
				return
			}
			ctx := contextutils.SetLogin(r.Context(), claims.Login)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getTokenStringFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("authorization header is missing or invalid")
	}
	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

func handleUnauthorizedError(w http.ResponseWriter) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}
