package jwt

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}
