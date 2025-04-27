package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

func GenerateToken(login, jwtSecret string, jwtExp time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"login": login,
		"exp":   time.Now().Add(jwtExp).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		logger.Logger.Errorf("Error signing token: %v", err)
		return "", err
	}
	return signedToken, nil
}
