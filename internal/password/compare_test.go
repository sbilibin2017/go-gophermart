package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestCompare(t *testing.T) {
	t.Run("successful comparison", func(t *testing.T) {
		password := "mySecurePassword"
		hashedPassword, err := Hash(password)
		require.NoError(t, err)
		err = Compare(hashedPassword, password)
		require.NoError(t, err)
	})

	t.Run("incorrect password", func(t *testing.T) {
		password := "mySecurePassword"
		incorrectPassword := "wrongPassword"
		hashedPassword, err := Hash(password)
		require.NoError(t, err)
		err = Compare(hashedPassword, incorrectPassword)
		require.Error(t, err)
		assert.Equal(t, err.Error(), bcrypt.ErrMismatchedHashAndPassword.Error())
	})

	t.Run("empty password", func(t *testing.T) {
		password := ""
		hashedPassword, err := Hash(password)
		require.NoError(t, err)
		err = Compare(hashedPassword, password)
		require.NoError(t, err)
	})

	t.Run("empty hashed password", func(t *testing.T) {
		password := "nonEmptyPassword"
		emptyHashedPassword := ""
		err := Compare(emptyHashedPassword, password)
		require.Error(t, err)
	})
}
