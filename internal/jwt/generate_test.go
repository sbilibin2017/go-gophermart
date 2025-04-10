package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	t.Run("successful token generation", func(t *testing.T) {
		secret := "mySecretKey"
		login := "user123"
		expireTime := time.Hour
		token, err := Generate(login, secret, expireTime)
		require.NoError(t, err)
		require.NotEmpty(t, token)
	})

	t.Run("token with short expiration", func(t *testing.T) {
		secret := "mySecretKey"
		login := "user123"
		expireTime := time.Second
		token, err := Generate(login, secret, expireTime)
		require.NoError(t, err)
		require.NotEmpty(t, token)
	})
}
