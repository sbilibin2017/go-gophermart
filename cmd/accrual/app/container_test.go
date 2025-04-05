package app

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestContainer_NewContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockConfig := NewMockContainerConfig(ctrl)
	container, err := NewContainer(mockConfig)
	assert.NoError(t, err)
	assert.NotNil(t, container)
	assert.Equal(t, mockConfig, container.Config)
	assert.NotNil(t, container.AccrualRouter)
}

func TestContainer_ConfigMethods(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockConfig := NewMockContainerConfig(ctrl)
	mockConfig.EXPECT().GetDatabaseURI().Return("postgres://user:password@localhost:5432/db").Times(1)
	mockConfig.EXPECT().GetRunAddress().Return(":8080").Times(1)
	assert.Equal(t, "postgres://user:password@localhost:5432/db", mockConfig.GetDatabaseURI())
	assert.Equal(t, ":8080", mockConfig.GetRunAddress())
}
