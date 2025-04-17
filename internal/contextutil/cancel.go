package contextutil

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

func NewCancelContext() (context.Context, context.CancelFunc) {
	logger.Logger.Info("Создание контекста для обработки сигналов завершения...")
	return signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
}
