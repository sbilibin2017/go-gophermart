package jwt

import (
	"github.com/golang-jwt/jwt/v4"
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
				return nil, ErrUnexpectedSigningMethod
			}
			return []byte(d.g.GetSecretKey()), nil
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
