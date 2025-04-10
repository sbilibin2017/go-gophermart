package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func Generate(login string, secret string, expireTime time.Duration) (string, error) {
	claims := &Claims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "test",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireTime)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
