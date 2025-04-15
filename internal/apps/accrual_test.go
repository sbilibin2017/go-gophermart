package apps

import (
	"testing"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAccrualApp(t *testing.T) {
	server, err := NewAccrualApp(configs.NewAccrualConfig("test", "test"), nil)
	require.NoError(t, err)
	assert.NotNil(t, server)
}
