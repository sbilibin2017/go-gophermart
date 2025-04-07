package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	secretKey := "my_secret_key"
	exp := time.Hour
	config := &configs.JWTConfig{
		SecretKey: secretKey,
		Exp:       exp,
	}
	login := "testuser"
	tokenString, err := Generate(config, login)
	require.NoError(t, err, "Generate should not return an error")
	assert.NotEmpty(t, tokenString, "Token string should not be empty")
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &Claims{})
	require.NoError(t, err, "Failed to parse token")
	claims, ok := token.Claims.(*Claims)
	require.True(t, ok, "Token claims should be of type *Claims")
	assert.Equal(t, login, claims.Login, "Token should have the correct Login claim")
	expectedExp := time.Now().Add(exp).Truncate(time.Second)
	actualExp := claims.ExpiresAt.Time.Truncate(time.Second)
	assert.Equal(t, expectedExp, actualExp, "Token expiration time should match the expected value")
}
