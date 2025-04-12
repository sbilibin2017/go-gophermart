package app

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewContainer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfig := NewMockContainerConfig(ctrl)

	mockConfig.EXPECT().GetRunAddress().Return(":8080").Times(1)
	mockConfig.EXPECT().GetJWTSecretKey().Return("secretkey").Times(0)
	mockConfig.EXPECT().GetJWTExpireTime().Return(24 * time.Hour).Times(0)

	container, err := NewContainer(mockConfig, nil)

	assert.NoError(t, err)
	assert.NotNil(t, container)
	assert.Equal(t, ":8080", container.Server.Addr)
}
