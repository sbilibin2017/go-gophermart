package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sbilibin2017/go-gophermart/internal/types"
	"github.com/stretchr/testify/assert"
)

type MockEncoderConfig struct {
	secretKey string
	exp       time.Duration
}

func (m *MockEncoderConfig) GetSecretKey() string {
	return m.secretKey
}

func (m *MockEncoderConfig) GetExp() time.Duration {
	return m.exp
}

func TestEncoder_Encode(t *testing.T) {
	secretKey := "mySecretKey"
	expiration := 24 * time.Hour
	claims := types.NewClaims(1)
	config := &MockEncoderConfig{
		secretKey: secretKey,
		exp:       expiration,
	}
	encoder := NewEncoder(config)
	tokenString, err := encoder.Encode(claims)
	assert.NoError(t, err, "Expected no error while encoding token")
	assert.NotEmpty(t, tokenString, "The generated token should not be empty")
	parsedToken, parseErr := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenIsNotValid
		}
		return []byte(secretKey), nil
	})
	assert.NoError(t, parseErr, "Expected no error while parsing the token")
	assert.NotNil(t, parsedToken, "Parsed token should not be nil")
	assert.True(t, parsedToken.Valid, "Token should be valid")
	parsedClaims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok, "Claims should be of type MapClaims")
	assert.Equal(t, float64(1), parsedClaims["user_id"], "UserID claim should be '1'")
	expTime := int64(parsedClaims["exp"].(float64))
	assert.Greater(t, expTime, time.Now().Unix(), "Expiration time should be in the future")
}
