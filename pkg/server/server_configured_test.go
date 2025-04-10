package server

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewServerConfigured(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAddr := NewMockServerConfiguredAddresser(ctrl)
	expectedAddr := ":9090"
	mockAddr.EXPECT().GetRunAddress().Return(expectedAddr)
	svrWithRouter := NewServerConfigured(mockAddr)
	assert.NotNil(t, svrWithRouter)
}
