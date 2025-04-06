package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type JWTDecoderConfig interface {
	GetSecretKey() string
}

type Decoder struct {
	c JWTDecoderConfig
}

func NewDecoder(c JWTDecoderConfig) *Decoder {
	return &Decoder{c: c}
}

var (
	ErrTokenIsNotValid = errors.New("token is not valid")
)

func (d *Decoder) Decode(tokenString string) (*types.Claims, error) {
	keyFunc := func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenIsNotValid
		}
		return []byte(d.c.GetSecretKey()), nil
	}
	claims := &types.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return nil, ErrTokenIsNotValid
	}
	if !token.Valid {
		return nil, ErrTokenIsNotValid
	}
	return claims, nil
}
