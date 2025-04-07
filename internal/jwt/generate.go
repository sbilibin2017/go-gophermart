package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
)

func Generate(config *configs.JWTConfig, login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Exp)),
		},
		Login: login,
	})
	tokenString, _ := token.SignedString([]byte(config.SecretKey))
	return tokenString, nil
}
