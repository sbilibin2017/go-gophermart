package contextutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCancelContext(t *testing.T) {
	ctx, cancel := NewCancelContext()
	assert.NotNil(t, ctx)
	assert.NotNil(t, cancel)

	select {
	case <-ctx.Done():
		t.Fatal("контекст не должен быть завершён до вызова cancel")
	default:
	}

	cancel()

	select {
	case <-ctx.Done():
	case <-time.After(1 * time.Second):
		t.Fatal("контекст должен быть завершён после вызова cancel")
	}
}
