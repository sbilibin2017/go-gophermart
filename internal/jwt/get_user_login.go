package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
)

func GetUserLogin(config *configs.JWTConfig, tokenString string) (string, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return []byte(config.SecretKey), nil
	})
	if err != nil {
		return "", err
	}
	return claims.Login, nil
}
