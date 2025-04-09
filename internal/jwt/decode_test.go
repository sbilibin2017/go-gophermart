package jwt

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDecode_ValidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockConfig := NewMockDecodeConfig(ctrl)
	secretKey := "test-secret-key"
	mockConfig.EXPECT().GetJWTSecretKey().Return(secretKey).Times(1)
	claims := jwt.MapClaims{
		"sub": "user123",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Fatalf("Failed to create token string: %v", err)
	}
	decodedClaims, err := Decode(tokenString, mockConfig)
	assert.NoError(t, err)
	assert.NotNil(t, decodedClaims)
	assert.IsType(t, &Claims{}, decodedClaims)
}
