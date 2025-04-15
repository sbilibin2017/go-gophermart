package runners

import (
	"context"

	"github.com/sbilibin2017/go-gophermart/internal/log"
)

type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

func RunServer(ctx context.Context, srv Server) error {
	log.Info("Запуск сервера...")

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Info("Контекст завершён, shutting down сервер...")
		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Error("Ошибка при завершении работы сервера", "ошибка", err)
		} else {
			log.Info("Сервер завершил работу корректно")
		}
		return err
	case err := <-errCh:
		if err != nil {
			log.Error("Ошибка на сервере", "ошибка", err)
		}
		return err
	}
}
