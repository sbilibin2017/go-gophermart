package server

import (
	"context"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

func Run(ctx context.Context, server Server) error {
	errChan := make(chan error)

	logger.Logger.Info("Сервер запускается...")

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		logger.Logger.Info("Получен сигнал о завершении работы сервера.")

		if err := server.Shutdown(ctx); err != nil {
			logger.Logger.Error("Ошибка при завершении сервера:", err)
			return err
		}

		logger.Logger.Info("Сервер успешно завершил работу.")
	case err := <-errChan:
		logger.Logger.Error("Ошибка сервера:", err)
		return err
	}

	return nil
}
