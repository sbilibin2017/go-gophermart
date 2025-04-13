package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMainFunction(t *testing.T) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		main()
	}()
	select {
	case <-done:
		assert.True(t, true)
	case <-time.After(5 * time.Second):
		t.Fatal("Main function did not finish in time")
	}
}
