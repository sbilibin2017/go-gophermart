package server

import (
	"context"

	"net/http"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/log"
)

type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

func Run(ctx context.Context, srv Server) error {
	go func() error {
		log.Logger.Infof("Запуск сервера...")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Logger.Fatalf("Ошибка сервера: %v", err)
			return err
		}
		return nil
	}()

	<-ctx.Done()
	log.Logger.Info("Ожидаем завершения работы...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Logger.Fatalf("Не удалось завершить работу сервера: %v", err)
		return err
	}

	log.Logger.Info("Сервер завершил работу корректно")

	return nil
}
