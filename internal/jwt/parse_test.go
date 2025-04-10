package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Run("successful parse", func(t *testing.T) {
		secret := "mySecretKey"
		login := "user123"
		expireTime := time.Hour

		token, err := Generate(login, secret, expireTime)
		require.NoError(t, err)

		parsedClaims, err := Parse(token, secret)

		require.NoError(t, err)
		assert.Equal(t, login, parsedClaims.Login)
		assert.NotNil(t, parsedClaims.ExpiresAt)
	})

	t.Run("expired token", func(t *testing.T) {
		secret := "mySecretKey"
		login := "user123"
		expireTime := -time.Hour

		token, err := Generate(login, secret, expireTime)
		require.NoError(t, err)

		parsedClaims, err := Parse(token, secret)

		assert.Equal(t, ErrExpiredToken, err)
		assert.Nil(t, parsedClaims)
	})

	t.Run("invalid secret", func(t *testing.T) {
		secret := "mySecretKey"
		invalidSecret := "wrongSecret"
		login := "user123"
		expireTime := time.Hour

		token, err := Generate(login, secret, expireTime)
		require.NoError(t, err)

		parsedClaims, err := Parse(token, invalidSecret)

		assert.Equal(t, ErrInvalidToken, err)
		assert.Nil(t, parsedClaims)
	})

	t.Run("invalid token format", func(t *testing.T) {
		invalidToken := "invalidTokenString"
		secret := "mySecretKey"

		parsedClaims, err := Parse(invalidToken, secret)

		assert.Equal(t, ErrInvalidToken, err)
		assert.Nil(t, parsedClaims)
	})
}
