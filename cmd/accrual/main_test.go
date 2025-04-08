package main

import (
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	go func() {
		main()
	}()
	time.Sleep(1 * time.Second)
}
