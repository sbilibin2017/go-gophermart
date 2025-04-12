package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTGeneratorConfigurer interface {
	GetJWTSecretKey() string
	GetJWTExpireTime() time.Duration
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(g.c.GetJWTExpireTime())),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(g.c.GetJWTSecretKey()))
	return tokenString
}
