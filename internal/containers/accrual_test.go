package containers

import (
	"testing"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/stretchr/testify/assert"
)

func TestNewAccrualContainer(t *testing.T) {
	mockConfig := &configs.JWTConfig{}
	container, err := NewAccrualContainer(mockConfig)
	assert.NoError(t, err)
	assert.NotNil(t, container)
	assert.NotNil(t, container.AccrualRouter)
}
