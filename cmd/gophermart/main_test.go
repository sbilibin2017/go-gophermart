package main

import (
	"context"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		main()
	}()

	time.Sleep(2 * time.Second)
	cancel()

	<-ctx.Done()
}
