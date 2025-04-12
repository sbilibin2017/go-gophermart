package server_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/server"
	"github.com/stretchr/testify/assert"
)

func TestRun_ServerStartsAndStopsGracefully(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := server.NewMockServer(ctrl)

	// Контекст с отменой, чтобы симулировать остановку сервера
	ctx, cancel := context.WithCancel(context.Background())

	// Настраиваем ожидания
	mockServer.EXPECT().ListenAndServe().Return(nil)
	mockServer.EXPECT().Shutdown(gomock.Any()).Return(nil)

	// Останавливаем сервер через 100 мс
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	err := server.Run(ctx, mockServer)
	assert.NoError(t, err)
}

func TestRun_ServerStartReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := server.NewMockServer(ctrl)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startErr := errors.New("start failed")

	// Возвращает ошибку при запуске, но не http.ErrServerClosed
	mockServer.EXPECT().ListenAndServe().Return(startErr)
	mockServer.EXPECT().Shutdown(gomock.Any()).Return(nil)

	// Останавливаем сервер через 100 мс
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	err := server.Run(ctx, mockServer)
	assert.NoError(t, err, "Shutdown should still succeed even if ListenAndServe fails")
}

func TestRun_ShutdownReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := server.NewMockServer(ctrl)

	ctx, cancel := context.WithCancel(context.Background())

	mockServer.EXPECT().ListenAndServe().Return(nil)
	mockServer.EXPECT().Shutdown(gomock.Any()).Return(errors.New("shutdown error"))

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	err := server.Run(ctx, mockServer)
	assert.EqualError(t, err, "shutdown error")
}
