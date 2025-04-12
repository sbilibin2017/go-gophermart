package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestCompare_Success(t *testing.T) {
	password := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)
	c := NewComparer()
	err = c.Compare(string(hashedPassword), password)
	assert.NoError(t, err)
}

func TestCompare_Failure(t *testing.T) {
	password := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)
	wrongPassword := "wrongpassword"
	c := NewComparer()
	err = c.Compare(string(hashedPassword), wrongPassword)
	assert.Error(t, err)
}

func TestCompare_EmptyPassword(t *testing.T) {
	password := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)
	emptyPassword := ""
	c := NewComparer()
	err = c.Compare(string(hashedPassword), emptyPassword)
	assert.Error(t, err)
}

func TestCompare_EmptyHashedPassword(t *testing.T) {
	password := "password123"
	emptyHashedPassword := ""
	c := NewComparer()
	err := c.Compare(string(emptyHashedPassword), password)
	assert.Error(t, err)
}
