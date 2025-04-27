package contexts

import (
	"context"
	"os/signal"
	"syscall"
)

func NewRunContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
}
