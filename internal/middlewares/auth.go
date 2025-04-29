package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/jwt"
)

func AuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := getTokenString(r)
			if err != nil {
				handleUnauthorizedError(w)
				return
			}
			login, err := jwt.GetLogin(tokenString, secretKey)
			if err != nil {
				handleUnauthorizedError(w)
				return
			}
			ctx := contextutils.SetLogin(r.Context(), login)
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
