package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/log"
)

type contextKey string

const userLoginKey contextKey = "userLogin"

func AuthMiddleware(config *configs.GophermartConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := getTokenHeader(r)
			if err != nil {
				log.Error("Failed to get token header", "error", err)
				http.Error(w, "Invalid token header", http.StatusUnauthorized)
				return
			}
			login, err := getUserLogin(config, tokenString)
			if err != nil {
				log.Error("Failed to parse or validate token", "error", err)
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}
			log.Info("Authenticated request", "login", login)
			r = withUserLoginContext(r, login)
			next.ServeHTTP(w, r)
		})
	}
}

func getTokenHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Error("Authorization header missing", "method", r.Method, "url", r.URL.Path)
		return "", errors.New("authorization header missing")
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		log.Error("Invalid authorization format", "method", r.Method, "url", r.URL.Path, "header", authHeader)
		return "", errors.New("invalid authorization format")
	}
	log.Info("Authorization header received", "method", r.Method, "url", r.URL.Path)
	return parts[1], nil
}

func getUserLogin(config *configs.GophermartConfig, tokenString string) (string, error) {
	log.Info("Parsing JWT token", "token", tokenString)
	claims := struct {
		jwt.RegisteredClaims
		Login string `json:"login"`
	}{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		return []byte(config.JWTSecretKey), nil
	})
	if err != nil {
		log.Error("Failed to parse JWT token", "token", tokenString, "error", err)
		return "", err
	}
	log.Info("Token successfully parsed", "login", claims.Login)
	return claims.Login, nil
}

func withUserLoginContext(r *http.Request, login string) *http.Request {
	log.Info("Adding user login to context", "login", login, "method", r.Method, "url", r.URL.Path)
	ctx := context.WithValue(r.Context(), userLoginKey, login)
	return r.WithContext(ctx)
}
