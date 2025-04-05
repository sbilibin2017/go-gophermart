package context

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStartContext_CancelOnSignal(t *testing.T) {
	ctx, cancel := NewStartContext()
	defer cancel()
	err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	require.NoError(t, err, "failed to send signal")
	select {
	case <-ctx.Done():
		assert.Equal(t, context.Canceled, ctx.Err(), "context should be canceled on signal")
	case <-time.After(2 * time.Second):
		t.Fatal("context was not canceled after signal")
	}
}
