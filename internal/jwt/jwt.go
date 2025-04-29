package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

type claims struct {
	jwt.RegisteredClaims
	Login string `json:"login"`
}

func GenerateTokenString(login string, secretKey string, exp time.Duration) (string, error) {
	claims := claims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "gophermart",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Преобразуем секретный ключ в []byte
	secretKeyBytes := []byte(secretKey)

	// Подписываем токен
	signedToken, err := token.SignedString(secretKeyBytes)
	if err != nil {
		logger.Logger.Errorw("Failed to sign token", "error", err, "login", login)
		return "", fmt.Errorf("failed to sign token: %v", err)
	}
	logger.Logger.Infow("Token successfully generated", "login", login, "token", signedToken)
	return signedToken, nil
}

func GetLogin(tokenString string, secretKey string) (string, error) {
	claims := &claims{}
	// Преобразуем секретный ключ в []byte
	secretKeyBytes := []byte(secretKey)

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (any, error) {
			return secretKeyBytes, nil
		})
	if err != nil {
		logger.Logger.Errorw("Failed to parse token", "error", err, "tokenString", tokenString)
		return "", fmt.Errorf("failed to parse token: %w", err)
	}
	if !token.Valid {
		logger.Logger.Errorw("Token is not valid", "tokenString", tokenString)
		return "", fmt.Errorf("token is not valid")
	}
	logger.Logger.Infow("Token parsed successfully", "login", claims.Login, "tokenString", tokenString)
	return claims.Login, nil
}
