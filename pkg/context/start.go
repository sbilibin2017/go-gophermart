package context

import (
	"context"
	"os/signal"
	"syscall"
)

func NewStartContext() (context.Context, context.CancelFunc) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	return ctx, cancel
}
