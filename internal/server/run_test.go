package server

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServer := NewMockServer(ctrl)
	mockServer.EXPECT().ListenAndServe().Return(nil).Times(1)
	mockServer.EXPECT().Shutdown(gomock.Any()).Return(nil).Times(1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		err := Run(ctx, mockServer)
		require.NoError(t, err)
	}()
	cancel()
	time.Sleep(100 * time.Millisecond)
}

func TestRun_ShutdownError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServer := NewMockServer(ctrl)
	mockServer.EXPECT().ListenAndServe().Return(nil).Times(1)
	mockServer.EXPECT().Shutdown(gomock.Any()).Return(errors.New("shutdown failed")).Times(1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		err := Run(ctx, mockServer)
		require.Error(t, err)
		assert.Equal(t, "shutdown failed", err.Error())
	}()
	cancel()
	time.Sleep(100 * time.Millisecond)
}

func TestRun_ListenAndServeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServer := NewMockServer(ctrl)
	mockServer.EXPECT().ListenAndServe().Return(errors.New("server failed")).Times(1)
	mockServer.EXPECT().Shutdown(gomock.Any()).Return(nil).Times(1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		err := Run(ctx, mockServer)
		require.Error(t, err)
		assert.Equal(t, "server failed", err.Error())
	}()
	cancel()
	time.Sleep(100 * time.Millisecond)
}
