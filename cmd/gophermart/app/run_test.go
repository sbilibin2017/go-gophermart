package app

import (
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	go func() {
		Run()
	}()
	time.Sleep(1 * time.Second)
}
