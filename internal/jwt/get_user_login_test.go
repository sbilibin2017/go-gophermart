package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
)

func GenerateTestToken(secretKey string, login string) (string, error) {
	claims := Claims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "test",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func TestGetUserLogin(t *testing.T) {
	secretKey := "mysecret"
	login := "testuser"
	tokenString, err := GenerateTestToken(secretKey, login)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	config := &configs.JWTConfig{
		SecretKey: secretKey,
	}
	result, err := GetUserLogin(config, tokenString)
	if err != nil {
		t.Fatalf("Error getting user login: %v", err)
	}
	if result != login {
		t.Errorf("Expected login %v, got %v", login, result)
	}
}

func TestGetUserLoginInvalidToken(t *testing.T) {
	secretKey := "mysecret"
	invalidToken := "invalid.token.string"
	config := &configs.JWTConfig{
		SecretKey: secretKey,
	}
	result, err := GetUserLogin(config, invalidToken)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if result != "" {
		t.Errorf("Expected empty result, got %v", result)
	}
}
