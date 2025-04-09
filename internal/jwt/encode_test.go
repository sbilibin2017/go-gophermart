package jwt

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockConfig := NewMockEncodeConfig(ctrl)
	secretKey := "mySecretKey"
	expiration := 1 * time.Hour
	login := "testuser"
	mockConfig.EXPECT().GetJWTSecretKey().Return(secretKey).Times(1)
	mockConfig.EXPECT().GetJWTExp().Return(expiration).Times(1)
	token, err := Encode(mockConfig, login)
	assert.NoError(t, err)
	assert.True(t, len(token) > 0)
	assert.True(t, token[:3] == "eyJ")
}
