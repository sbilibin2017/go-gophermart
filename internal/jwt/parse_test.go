package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func generateTestToken(secret string, claims Claims, method jwt.SigningMethod) (string, error) {
	token := jwt.NewWithClaims(method, claims)
	return token.SignedString([]byte(secret))
}

func TestJWTParser_Parse_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockConfig := NewMockJWTParserConfigurer(ctrl)
	parser := &JWTParser{c: mockConfig}
	secret := "testsecret"
	mockConfig.EXPECT().GetSecretKey().Return(secret)
	expectedClaims := Claims{
		Login: "12345",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	tokenStr, err := generateTestToken(secret, expectedClaims, jwt.SigningMethodHS256)
	assert.NoError(t, err)
	parsedClaims, err := parser.Parse(tokenStr)
	assert.NoError(t, err)
	assert.Equal(t, expectedClaims.Login, parsedClaims.Login)
}

func TestJWTParser_Parse_InvalidSignatureAlgorithm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockConfig := NewMockJWTParserConfigurer(ctrl)
	parser := &JWTParser{c: mockConfig}
	secret := "testsecret"
	mockConfig.EXPECT().GetSecretKey().Return(secret).Times(0)
	claims := Claims{
		Login: "54321",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	tokenStr, err := generateTestToken(secret, claims, jwt.SigningMethodHS384)
	assert.NoError(t, err)
	parsedClaims, err := parser.Parse(tokenStr)
	assert.Nil(t, parsedClaims)
	assert.ErrorIs(t, err, ErrInvalidToken)
}
