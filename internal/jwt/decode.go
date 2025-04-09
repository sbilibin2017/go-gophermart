package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrTokenIsNotParsed        = errors.New("token is not parsed")
	ErrInvalidTokenClaims      = errors.New("invalid token claims")
)

type DecodeConfig interface {
	GetJWTSecretKey() string
}

func Decode(tokenStr string, config DecodeConfig) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		return []byte(config.GetJWTSecretKey()), nil
	})
	if err != nil {
		return nil, ErrTokenIsNotParsed
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidTokenClaims
	}
	return claims, nil
}
