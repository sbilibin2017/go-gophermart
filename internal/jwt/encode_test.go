package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	login := "testuser"
	token, err := Encode("secretkey", time.Hour*24, login)
	assert.NoError(t, err, "Token should be signed without error")
	assert.NotEmpty(t, token, "Token should not be empty")
	assert.Contains(t, token, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", "Token should start with the expected JWT header")
}
