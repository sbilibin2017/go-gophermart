package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain_Run(t *testing.T) {
	originalExitFunc := exitFunc
	defer func() { exitFunc = originalExitFunc }()
	exitCode := make(chan int)
	exitFunc = func(code int) {
		exitCode <- code
	}
	go main()
	code := <-exitCode
	assert.Equal(t, 0, code, "The exit code should be 0")
}
