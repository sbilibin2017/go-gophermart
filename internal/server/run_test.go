package server_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/server"
)

func TestRunWithGracefulShutdown(t *testing.T) {
	logger.Init()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := server.NewMockServer(ctrl)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockServer.EXPECT().ListenAndServe().Return(errors.New("server error")).Times(1)
	mockServer.EXPECT().Shutdown(gomock.Any()).Return(nil).Times(1)

	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	err := server.RunWithGracefulShutdown(ctx, mockServer)

	assert.NoError(t, err)
}

func TestRunWithGracefulShutdown_ServerListenAndServeSucceeds_ShutdownFails(t *testing.T) {
	logger.Init()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := server.NewMockServer(ctrl)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockServer.EXPECT().ListenAndServe().Return(nil).Times(1)
	mockServer.EXPECT().Shutdown(gomock.Any()).Return(errors.New("shutdown failed")).Times(1)

	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	err := server.RunWithGracefulShutdown(ctx, mockServer)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "shutdown failed")
}
