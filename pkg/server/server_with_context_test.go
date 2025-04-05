package server

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestServer_Start_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHTTPServer := NewMockHTTPServer(ctrl)
	mockHTTPServer.
		EXPECT().
		ListenAndServe().
		DoAndReturn(func() error {
			time.Sleep(100 * time.Millisecond)
			return http.ErrServerClosed
		})
	mockHTTPServer.
		EXPECT().
		Shutdown(gomock.Any()).
		Return(nil)
	s := NewServerWithContext(mockHTTPServer)
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()
	err := s.Start(ctx)
	assert.NoError(t, err)
}

func TestServer_Start_ListenAndServe_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	expectedErr := errors.New("ListenAndServe failed")
	mockHTTPServer := NewMockHTTPServer(ctrl)
	mockHTTPServer.
		EXPECT().
		ListenAndServe().
		Return(expectedErr)
	s := NewServerWithContext(mockHTTPServer)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := s.Start(ctx)
	assert.Equal(t, expectedErr, err)
}

func TestServer_Start_Shutdown_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHTTPServer := NewMockHTTPServer(ctrl)
	mockHTTPServer.
		EXPECT().
		ListenAndServe().
		DoAndReturn(func() error {
			time.Sleep(100 * time.Millisecond)
			return http.ErrServerClosed
		})
	shutdownErr := errors.New("shutdown failed")
	mockHTTPServer.
		EXPECT().
		Shutdown(gomock.Any()).
		Return(shutdownErr)
	s := NewServerWithContext(mockHTTPServer)
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()
	err := s.Start(ctx)
	assert.Equal(t, shutdownErr, err)
}

func TestServer_Start_ContextCancel_ShutdownSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHTTPServer := NewMockHTTPServer(ctrl)
	mockHTTPServer.
		EXPECT().
		ListenAndServe().
		DoAndReturn(func() error {
			time.Sleep(100 * time.Millisecond)
			return http.ErrServerClosed
		})
	mockHTTPServer.
		EXPECT().
		Shutdown(gomock.Any()).
		Return(nil)
	s := NewServerWithContext(mockHTTPServer)
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()
	err := s.Start(ctx)
	assert.NoError(t, err)
}
