package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTGeneratorConfigurer interface {
	GetSecretKey() string
	GetExpireTime() time.Duration
}

type JWTGenerator struct {
	c JWTGeneratorConfigurer
}

func NewJWTGenerator(c JWTGeneratorConfigurer) *JWTGenerator {
	return &JWTGenerator{c: c}
}

func (g *JWTGenerator) Generate(login string) string {
	claims := &Claims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(g.c.GetExpireTime())),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(g.c.GetSecretKey()))
	return tokenString
}
