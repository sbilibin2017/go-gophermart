package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	claims := &Claims{
		Login: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("secretkey"))
	assert.NoError(t, err, "Failed to sign token")
	decodedClaims, err := Decode(signedToken, "secretkey")
	assert.NoError(t, err, "Failed to decode token")
	assert.Equal(t, claims.Login, decodedClaims.Login)
	assert.NotNil(t, decodedClaims.RegisteredClaims)
	assert.True(t, decodedClaims.RegisteredClaims.ExpiresAt.Time.After(time.Now()))
}
