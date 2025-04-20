package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain_CallsFlagsAndRun(t *testing.T) {
	calledFlags := false
	calledRun := false

	originalFlags := flagsFunc
	originalRun := runFunc

	defer func() {
		flagsFunc = originalFlags
		runFunc = originalRun
	}()

	flagsFunc = func() {
		calledFlags = true
	}
	runFunc = func() {
		calledRun = true
	}

	main()

	assert.True(t, calledFlags, "flagsFunc should be called from main()")
	assert.True(t, calledRun, "runFunc should be called from main()")
}
