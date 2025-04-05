package context

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"context"
)

func TestNewShutdownContext(t *testing.T) {
	ctx, cancel := NewShutdownContext()
	require.NotNil(t, ctx)
	require.NotNil(t, cancel)
	select {
	case <-ctx.Done():
		t.Fatal("context should not be canceled yet")
	default:
	}
	cancel()
	select {
	case <-ctx.Done():
		assert.Equal(t, context.Canceled, ctx.Err(), "context should be canceled")
	case <-time.After(100 * time.Millisecond):
		t.Fatal("context should be canceled after calling cancel")
	}
}

func TestNewShutdownContext_Timeout(t *testing.T) {
	ctx, _ := NewShutdownContext()
	select {
	case <-ctx.Done():
		t.Fatal("context canceled too early")
	case <-time.After(1 * time.Second):
	}
	select {
	case <-ctx.Done():
		assert.Equal(t, context.DeadlineExceeded, ctx.Err(), "context should timeout after 5 seconds")
	case <-time.After(6 * time.Second):
		t.Fatal("context should have timed out")
	}
}
