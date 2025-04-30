package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClaims(t *testing.T) {

	user := &User{Login: "testuser"}
	issuer := "testissuer"
	expiration := time.Hour * 24

	currentTime := time.Now()

	claims := NewClaims(user, issuer, expiration)

	assert.Equal(t, user.Login, claims.Login, "Login in claims should match the user login")
	assert.Equal(t, issuer, claims.Issuer, "Issuer should match the provided issuer")
	assert.WithinDuration(t, currentTime.Add(expiration), claims.ExpiresAt.Time, time.Second, "Expiration time should be correct")
	assert.WithinDuration(t, currentTime, claims.IssuedAt.Time, time.Second, "IssuedAt time should be correct")
	assert.NotEmpty(t, claims.Issuer, "Issuer should not be empty")
	assert.NotNil(t, claims.ExpiresAt, "ExpiresAt should not be nil")
	assert.NotNil(t, claims.IssuedAt, "IssuedAt should not be nil")
}
