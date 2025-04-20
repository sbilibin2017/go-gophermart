package contextutils

import (
	"testing"
)

func TestNewCancelContext(t *testing.T) {
	_, cancel := NewCancelContext()
	defer cancel()
}
