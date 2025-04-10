package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestJWTDecoder_Decode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGetter := NewMockJWTSecretKeyGetter(ctrl)
	secretKey := "test_secret"
	mockGetter.EXPECT().GetSecretKey().Return(secretKey).Times(1)
	claims := &Claims{
		Login: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secretKey))
	assert.NoError(t, err)
	decoder := JWTDecoder{g: mockGetter}
	decodedClaims, err := decoder.Decode(tokenStr)
	assert.NoError(t, err)
	assert.NotNil(t, decodedClaims)
	assert.Equal(t, "testuser", decodedClaims.Login)
}
