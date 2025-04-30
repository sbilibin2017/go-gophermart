package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestGenerateTokenString_Success(t *testing.T) {
	user := &types.User{Login: "testuser"}
	config := &configs.JWTConfig{
		SecretKey: "secret",
		Issuer:    "testissuer",
		Exp:       time.Hour,
	}
	tokenString, err := GenerateTokenString(config, user)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)
}

func TestGetClaims_Success(t *testing.T) {
	user := &types.User{Login: "testuser"}
	config := &configs.JWTConfig{
		SecretKey: "secret",
		Issuer:    "testissuer",
		Exp:       time.Hour,
	}
	tokenString, err := GenerateTokenString(config, user)
	assert.NoError(t, err)
	claims, err := DecodeTokenString(tokenString, config)
	assert.NoError(t, err)
	assert.Equal(t, user.Login, claims.Login)
	assert.Equal(t, config.Issuer, claims.Issuer)
	assert.WithinDuration(t, time.Now().Add(config.Exp), claims.ExpiresAt.Time, time.Second)
}

func TestGetClaims_InvalidToken2(t *testing.T) {
	config := &configs.JWTConfig{
		SecretKey: "secret",
		Issuer:    "testissuer",
		Exp:       time.Hour,
	}
	invalidToken := "invalid.token.string"
	claims, err := DecodeTokenString(invalidToken, config)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestGetClaims_ExpiredToken(t *testing.T) {
	user := &types.User{Login: "testuser"}
	config := &configs.JWTConfig{
		SecretKey: "secret",
		Issuer:    "testissuer",
		Exp:       -time.Hour,
	}
	tokenString, err := GenerateTokenString(config, user)
	assert.NoError(t, err)
	claims, err := DecodeTokenString(tokenString, config)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestGetClaims_InvalidToken(t *testing.T) {
	config := &configs.JWTConfig{
		SecretKey: "secret",
		Issuer:    "testissuer",
		Exp:       time.Hour,
	}
	user := &types.User{Login: "testuser"}
	claims := types.NewClaims(user, config.Issuer, config.Exp)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, *claims)
	signedToken, err := token.SignedString([]byte(config.SecretKey))
	assert.NoError(t, err)
	invalidToken := signedToken + "invalid"
	parsedClaims, err := DecodeTokenString(invalidToken, config)
	assert.Error(t, err)
	assert.Nil(t, parsedClaims)
}
