package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

func AuthMiddleware(
	config *configs.JWTConfig,
	tokenDecoder func(tokenString string, config *configs.JWTConfig) (*types.Claims, error),
	claimsContextSetter func(ctx context.Context, claims *types.Claims) context.Context,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := getTokenStringFromHeader(r)
			if err != nil {
				handleUnauthorizedError(w)
				return
			}
			claims, err := tokenDecoder(tokenString, config)
			if err != nil {
				handleUnauthorizedError(w)
				return
			}
			ctx := claimsContextSetter(r.Context(), claims)
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
