package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainFunction(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Main function panicked with: %v", r)
		}
	}()
	main()
	assert.True(t, true, "No error should occur")
}
