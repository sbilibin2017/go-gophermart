package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

func GenerateTokenString(
	config *configs.JWTConfig,
	user *types.User,
) (string, error) {
	claims := types.NewClaims(user, config.Issuer, config.Exp)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, *claims)
	signedToken, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}
	return signedToken, nil
}

func DecodeTokenString(
	tokenString string,
	config *configs.JWTConfig,
) (*types.Claims, error) {
	claims := &types.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (any, error) {
			return []byte(config.SecretKey), nil
		})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}
	return claims, nil
}
