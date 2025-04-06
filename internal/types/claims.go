package types

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	jwt.RegisteredClaims
	*User
}

func NewClaims(userID int64) *Claims {
	return &Claims{
		User: &User{
			UserID: &UserID{UserID: userID},
		},
	}
}
