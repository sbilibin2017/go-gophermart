package contextutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCancelContext(t *testing.T) {
	ctx, cancel := NewCancelContext()
	assert.NotNil(t, ctx, "context should not be nil")
	assert.NotNil(t, cancel, "cancel function should not be nil")

	select {
	case <-ctx.Done():
		t.Fatal("context should not be canceled yet")
	default:
	}

	cancel()

	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("context should be canceled after calling cancel")
	}
}
