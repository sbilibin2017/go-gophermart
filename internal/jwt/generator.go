package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTGenerator struct {
	SecretKey      string
	ExpirationTime time.Duration
}

func NewJWTGenerator(secretKey string, expirationTime time.Duration) *JWTGenerator {
	return &JWTGenerator{
		SecretKey:      secretKey,
		ExpirationTime: expirationTime,
	}
}

func (g *JWTGenerator) Generate(payload map[string]any) *string {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(g.ExpirationTime).Unix(),
	}
	for key, value := range payload {
		claims[key] = value
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(g.SecretKey))
	if err != nil {
		return nil
	}
	return &tokenString
}
