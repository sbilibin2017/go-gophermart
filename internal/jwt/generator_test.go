package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestJWTGenerator_Generate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockConfig := NewMockJWTGeneratorConfigurer(ctrl)
	secretKey := "testsecret"
	expireTime := time.Hour
	mockConfig.EXPECT().GetSecretKey().Return(secretKey).Times(1)
	mockConfig.EXPECT().GetExpireTime().Return(expireTime).Times(1)
	generator := NewJWTGenerator(mockConfig)
	login := "testuser"
	token := generator.Generate(login)
	assert.NotEmpty(t, token)
	parsedClaims := &Claims{}
	_, err := jwt.ParseWithClaims(token, parsedClaims, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
	assert.NoError(t, err)
	assert.Equal(t, login, parsedClaims.Login)
	assert.NotNil(t, parsedClaims.RegisteredClaims.ExpiresAt)
	assert.True(t, parsedClaims.RegisteredClaims.ExpiresAt.Time.Before(time.Now().Add(expireTime)))
}

func TestJWTGenerator_Generate_EmptySecretKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockConfig := NewMockJWTGeneratorConfigurer(ctrl)
	mockConfig.EXPECT().GetSecretKey().Return("").Times(1)
	mockConfig.EXPECT().GetExpireTime().Return(time.Hour).Times(1)
	generator := NewJWTGenerator(mockConfig)
	login := "testuser"
	token := generator.Generate(login)
	assert.NotEmpty(t, token)
	parsedClaims := &Claims{}
	_, err := jwt.ParseWithClaims(token, parsedClaims, func(token *jwt.Token) (any, error) {
		return []byte(""), nil 
	})
	assert.NoError(t, err)
}

func TestJWTGenerator_Generate_InvalidSecretKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockConfig := NewMockJWTGeneratorConfigurer(ctrl)
	secretKey := "testsecret"
	expireTime := time.Hour
	mockConfig.EXPECT().GetSecretKey().Return(secretKey).Times(1)
	mockConfig.EXPECT().GetExpireTime().Return(expireTime).Times(1)
	generator := NewJWTGenerator(mockConfig)
	login := "testuser"
	token := generator.Generate(login)
	assert.NotEmpty(t, token)
	parsedClaims := &Claims{}
	_, err := jwt.ParseWithClaims(token, parsedClaims, func(token *jwt.Token) (any, error) {
		return []byte("wrongsecret"), nil 
	})
	assert.Error(t, err)
}
