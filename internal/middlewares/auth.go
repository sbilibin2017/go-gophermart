package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"go.uber.org/zap"
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
				if logger.Logger != nil {
					logger.Logger.Warn("AuthMiddleware: Authorization header missing")
				}
				http.Error(w, ErrMissingAuthHeader, http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				if logger.Logger != nil {
					logger.Logger.Warn("AuthMiddleware: Invalid authorization header format")
				}
				http.Error(w, ErrInvalidAuthHeader, http.StatusUnauthorized)
				return
			}

			claims, err := parseToken(parts[1], jwtKey)
			if err != nil {
				if logger.Logger != nil {
					logger.Logger.Warn("AuthMiddleware: Invalid or expired token", zap.Error(err))
				}
				http.Error(w, ErrInvalidOrExpiredToken, http.StatusUnauthorized)
				return
			}

			if logger.Logger != nil {
				logger.Logger.Info("AuthMiddleware: Successfully authenticated user", zap.String("login", claims.Login))
			}

			ctx := context.WithValue(r.Context(), claimsContextKey, claims.Login)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetLogin(r *http.Request) (string, error) {
	login, ok := r.Context().Value(claimsContextKey).(string)
	if !ok {
		return "", errors.New("context error")

	}
	return login, nil
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
