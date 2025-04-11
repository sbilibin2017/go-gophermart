package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token is expired")
)

type JWTParserConfigurer interface {
	GetSecretKey() string
}

type JWTParser struct {
	c JWTParserConfigurer
}

func (p *JWTParser) Parse(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, _ := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidToken
		}
		return []byte(p.c.GetSecretKey()), nil
	})
	if !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
