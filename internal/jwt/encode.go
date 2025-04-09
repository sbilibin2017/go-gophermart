package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
)

var (
	ErrTokenIsNotSigned = errors.New("token is not signed")
)

func Encode(config *configs.GophermartConfig, login string) (string, error) {
	claims := &Claims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.JWTExp)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.JWTSecretKey))
	if err != nil {
		return "", ErrTokenIsNotSigned
	}
	return signedToken, nil
}
