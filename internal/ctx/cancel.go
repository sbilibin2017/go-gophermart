package ctx

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/sbilibin2017/go-gophermart/internal/log"
)

func NewCancelContext() (context.Context, context.CancelFunc) {
	log.Info("Создание нового контекста для отмены с уведомлениями о сигналах",
		"сигналы", []string{"SIGINT", "SIGTERM"},
	)
	return signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
}
