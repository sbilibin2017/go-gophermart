package server

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewServerConfigured(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAddr := "127.0.0.1:8000"
	mockAddresser := NewMockAddresser(ctrl)
	mockAddresser.EXPECT().GetRunAddress().Return(mockAddr)
	result := NewServerConfigured(mockAddresser)
	assert.NotNil(t, result)
	assert.NotNil(t, result.server)
	assert.NotNil(t, result.router)
}
