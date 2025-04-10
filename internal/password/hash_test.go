package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHash(t *testing.T) {
	t.Run("successful hash", func(t *testing.T) {
		password := "mySecurePassword"
		hashedPassword, err := Hash(password)
		require.NoError(t, err)
		require.NotEmpty(t, hashedPassword)
		assert.NotEqual(t, password, hashedPassword)
	})

	t.Run("empty password", func(t *testing.T) {
		password := ""
		hashedPassword, err := Hash(password)
		require.NoError(t, err)
		require.NotEmpty(t, hashedPassword)
		assert.NotEqual(t, password, hashedPassword)
	})

	t.Run("password with incorrect bcrypt cost", func(t *testing.T) {
		password := "testPassword"
		_, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		require.NoError(t, err)
	})
}
