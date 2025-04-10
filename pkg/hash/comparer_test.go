package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestComparer_Compare(t *testing.T) {
	password := "mysecurepassword"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err, "bcrypt hashing should not return an error")
	comparer := NewComarer()
	t.Run("matching passwords", func(t *testing.T) {
		err := comparer.Compare(string(hashedPassword), password)
		assert.NoError(t, err, "Compare should succeed when passwords match")
	})
	t.Run("non-matching passwords", func(t *testing.T) {
		err := comparer.Compare(string(hashedPassword), "wrongpassword")
		assert.Error(t, err, "Compare should return an error when passwords don't match")
	})
}
