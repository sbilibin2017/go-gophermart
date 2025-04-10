package jwt

import (
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestJWTEncoder_Encode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfig := NewMockJWTEncodeConfig(ctrl)

	secretKey := "test_secret"
	expDuration := time.Hour

	mockConfig.EXPECT().GetSecretKey().Return(secretKey).Times(1)
	mockConfig.EXPECT().GetExp().Return(expDuration).Times(1)

	encoder := JWTEncoder{c: mockConfig}

	token, err := encoder.Encode("testuser")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
