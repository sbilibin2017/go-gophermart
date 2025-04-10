package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHasher_Hash(t *testing.T) {
	hasher := NewHasher()
	password := "securepassword"
	hashedPassword, err := hasher.Hash(password)
	assert.NoError(t, err, "Hashing should not return an error")
	assert.NotEmpty(t, hashedPassword, "Hashed password should not be empty")
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	assert.NoError(t, err, "bcrypt should validate the password successfully")
}
