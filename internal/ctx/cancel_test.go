package ctx

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCancelContext(t *testing.T) {
	ctx, cancel := NewCancelContext()
	assert.NotNil(t, ctx)
	assert.Nil(t, ctx.Err())
	cancel()
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, context.Canceled, ctx.Err())
}
