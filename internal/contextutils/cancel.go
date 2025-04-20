package contextutils

import (
	"context"
	"os/signal"
	"syscall"
)

func NewCancelContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
}
