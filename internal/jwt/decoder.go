package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/sbilibin2017/go-gophermart/internal/errors"
)

type JWTSecretKeyGetter interface {
	GetSecretKey() string
}

type JWTDecoder struct {
	g JWTSecretKeyGetter
}

func (d JWTDecoder) Decode(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.ErrUnexpectedSigningMethod
			}
			return []byte(d.g.GetSecretKey()), nil
		})
	if err != nil {
		return nil, errors.ErrTokenIsNotParsed
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.ErrInvalidTokenClaims
	}
	return claims, nil
}
