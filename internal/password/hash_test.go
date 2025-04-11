package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHash_Success(t *testing.T) {
	password := "password123"
	hasher := NewHasher()
	hashedPassword := hasher.Hash(password)
	assert.NotEmpty(t, hashedPassword)
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	assert.NoError(t, err, "hashed password should match the original password")
}
