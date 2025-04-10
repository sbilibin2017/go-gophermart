package ctx

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCancelContext(t *testing.T) {
	ctx, cancel := NewCancelContext()
	require.NotNil(t, ctx)
	require.NotNil(t, cancel)
	defer cancel()
	process, err := os.FindProcess(os.Getpid())
	require.NoError(t, err)
	err = process.Signal(syscall.SIGINT)
	require.NoError(t, err)
	select {
	case <-ctx.Done():
		assert.Equal(t, context.Canceled, ctx.Err(), "context should be canceled")
	case <-time.After(1 * time.Second):
		t.Fatal("context was not canceled after sending SIGINT")
	}
}
