package server

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestRun_ServerStartAndShutdown(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSrv := NewMockServer(ctrl)
	mockSrv.EXPECT().ListenAndServe().DoAndReturn(func() error {
		time.Sleep(100 * time.Millisecond)
		return http.ErrServerClosed
	})
	mockSrv.EXPECT().Shutdown(gomock.Any()).Return(nil)
	ctx, cancel := context.WithCancel(context.Background())
	go Run(ctx, mockSrv)
	time.Sleep(50 * time.Millisecond)
	cancel()
	time.Sleep(200 * time.Millisecond)
}

func TestRun_ServerStartError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSrv := NewMockServer(ctrl)
	mockErr := errors.New("unexpected failure")
	mockSrv.EXPECT().ListenAndServe().Return(mockErr)
	mockSrv.EXPECT().Shutdown(gomock.Any()).Return(nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go Run(ctx, mockSrv)
	time.Sleep(200 * time.Millisecond)
	cancel()
	time.Sleep(200 * time.Millisecond)
}

func TestRun_ShutdownError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSrv := NewMockServer(ctrl)
	mockSrv.EXPECT().ListenAndServe().DoAndReturn(func() error {
		time.Sleep(50 * time.Millisecond)
		return http.ErrServerClosed
	})
	mockSrv.EXPECT().Shutdown(gomock.Any()).Return(errors.New("shutdown failed"))
	ctx, cancel := context.WithCancel(context.Background())
	go Run(ctx, mockSrv)
	time.Sleep(100 * time.Millisecond)
	cancel()
	time.Sleep(100 * time.Millisecond)
}
