package server

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerWithRouter_Start_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServer := NewMockHTTPServer(ctrl)
	router := chi.NewRouter()
	s := NewServerWithRouter(mockServer, router)
	mockServer.EXPECT().ListenAndServe().Return(nil)
	mockServer.EXPECT().Shutdown(gomock.Any()).Return(nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		err := s.Run(ctx)
		require.NoError(t, err)
	}()
	cancel()
	time.Sleep(1 * time.Second)
}

func TestStart_ServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHTTPServer := NewMockHTTPServer(ctrl)
	mockHTTPServer.EXPECT().ListenAndServe().Return(errors.New("server failed to start"))
	router := chi.NewRouter()
	s := NewServerWithRouter(mockHTTPServer, router)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := s.start(ctx)
	assert.Error(t, err, "server failed to start")
}

func TestServerWithRouter_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServer := NewMockHTTPServer(ctrl)
	router := chi.NewRouter()
	s := NewServerWithRouter(mockServer, router)
	mockServer.EXPECT().Shutdown(gomock.Any()).Return(nil)
	go func() {
		err := s.stop()
		assert.NoError(t, err)
	}()
	time.Sleep(1 * time.Second)
}

func TestServerWithRouter_Stop_WithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServer := NewMockHTTPServer(ctrl)
	router := chi.NewRouter()
	s := NewServerWithRouter(mockServer, router)
	mockServer.EXPECT().Shutdown(gomock.Any()).Return(errors.New("shutdown error"))
	go func() {
		err := s.stop()
		assert.EqualError(t, err, "shutdown error")
	}()

	time.Sleep(1 * time.Second)
}

func TestServerWithRouter_AddRouter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServer := NewMockHTTPServer(ctrl)
	mainRouter := chi.NewRouter()
	s := NewServerWithRouter(mockServer, mainRouter)
	newRouter := chi.NewRouter()
	mockServer.EXPECT().SetHandler(mainRouter)
	s.AddRouter(newRouter)
	assert.NotNil(t, mainRouter)
}
