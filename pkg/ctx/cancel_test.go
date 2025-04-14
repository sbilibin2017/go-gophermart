package ctx

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCancelContext(t *testing.T) {
	ctx, cancel := NewCancelContext()
	require.NotNil(t, ctx)
	require.NotNil(t, cancel)

	select {
	case <-ctx.Done():
		t.Fatal("context should not be canceled initially")
	default:
	}

	cancel()

	select {
	case <-ctx.Done():
		assert.Equal(t, context.Canceled, ctx.Err(), "context should be canceled")
	case <-time.After(time.Second):
		t.Fatal("context was not canceled in time")
	}
}
