package main

import (
	"errors"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain_Success(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	exitCalled := false
	exitFunc = func(code int) {
		exitCalled = (code == 0)
	}
	runFunc = func() error {
		return nil
	}
	go func() {
		main()
	}()
	time.Sleep(100 * time.Millisecond)
	assert.True(t, exitCalled)
}

func TestMain_Error(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	exitCalled := false
	exitFunc = func(code int) {
		exitCalled = (code == 0)
	}
	runFunc = func() error {
		return errors.New("err")
	}
	go func() {
		main()
	}()
	time.Sleep(100 * time.Millisecond)
	assert.True(t, exitCalled)
}
