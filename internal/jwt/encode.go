package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrTokenIsNotSigned = errors.New("token is not signed")
)

type EncodeConfig interface {
	GetJWTSecretKey() string
	GetJWTExp() time.Duration
}

func Encode(c EncodeConfig, login string) (string, error) {
	claims := &Claims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(c.GetJWTExp())),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(c.GetJWTSecretKey()))
	if err != nil {
		return "", ErrTokenIsNotSigned
	}
	return signedToken, nil
}
