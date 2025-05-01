package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(
	jwtSecretKey string,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := getTokenStringFromHeader(r)
			if err != nil {
				handleUnauthorizedError(w)
				return
			}
			login, err := getLoginFromToken(tokenString, jwtSecretKey)
			if err != nil {
				handleUnauthorizedError(w)
				return
			}
			ctx := setLoginToContext(r.Context(), login)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GenerateTokenString(
	jwtSecretKey string,
	jwtExp time.Duration,
	issuer string,
	login string,
) (string, error) {
	claims := struct {
		jwt.RegisteredClaims
		Login string `json:"login"`
	}{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}
	return signedToken, nil
}

func GetLoginFromContext(ctx context.Context) (string, error) {
	login, ok := ctx.Value(contextLoginKey{}).(string)
	if !ok {
		return "", errors.New("claims are not in context")
	}
	return login, nil
}

func getLoginFromToken(
	tokenString string,
	jwtSecretKey string,
) (string, error) {
	claims := &struct {
		jwt.RegisteredClaims
		Login string `json:"login"`
	}{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}
	if !token.Valid {
		return "", fmt.Errorf("token is not valid")
	}
	return claims.Login, nil
}

type contextLoginKey struct{}

func setLoginToContext(ctx context.Context, login string) context.Context {
	return context.WithValue(ctx, contextLoginKey{}, login)
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
