package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type JWTDecoder struct {
	SecretKey []byte
}

func NewJWTDecoder(secretKey []byte) *JWTDecoder {
	return &JWTDecoder{
		SecretKey: secretKey,
	}
}

func (d *JWTDecoder) Decode(tokenString string) (map[string]any, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return d.SecretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
