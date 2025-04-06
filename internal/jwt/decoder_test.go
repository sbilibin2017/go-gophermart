package jwt

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sbilibin2017/go-gophermart/internal/types"
	"github.com/stretchr/testify/assert"
)

type MockDecoderConfig struct {
	secretKey string
}

func (m *MockDecoderConfig) GetSecretKey() string {
	return m.secretKey
}

func TestDecoder_Decode(t *testing.T) {
	secretKey := "12345"
	claims := types.NewClaims(1)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}
	config := &MockDecoderConfig{secretKey: secretKey}
	decoder := NewDecoder(config)
	decodedClaims, err := decoder.Decode(signedToken)
	assert.NoError(t, err, "Expected no error while decoding token")
	assert.NotNil(t, decodedClaims, "Decoded claims should not be nil")
	assert.Equal(t, claims.UserID, decodedClaims.UserID, "UserID should be the same")
}
