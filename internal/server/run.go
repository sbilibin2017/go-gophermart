package server

import (
	"context"

	"github.com/sbilibin2017/go-gophermart/internal/log"
)

type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

func Run(ctx context.Context, server Server) error {
	errChan := make(chan error)

	// Start the server in a goroutine
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			errChan <- err // Send error to the channel
		}
	}()

	// Wait for either the context to be done or the server to return an error
	select {
	case <-ctx.Done():
		// If the context is done, shut down the server
		if err := server.Shutdown(ctx); err != nil {
			log.Logger.Error("Ошибка при завершении сервера:", err)
			return err
		}
	case err := <-errChan:
		// If an error occurs in the server, log it and return the error
		log.Logger.Error("Ошибка сервера:", err)
		return err
	}

	return nil
}
