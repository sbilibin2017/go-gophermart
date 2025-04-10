package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTEncodeConfig interface {
	GetSecretKey() string
	GetExp() time.Duration
}

type JWTEncoder struct {
	c JWTEncodeConfig
}

func (e JWTEncoder) Encode(login string) (string, error) {
	claims := &Claims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(e.c.GetExp())),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(e.c.GetSecretKey()))
	if err != nil {
		return "", ErrTokenIsNotSigned
	}
	return signedToken, nil
}
