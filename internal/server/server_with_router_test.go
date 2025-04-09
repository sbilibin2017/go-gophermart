package server

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerWithRouter_Start_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServer := NewMockHTTPServer(ctrl)
	router := chi.NewRouter()
	s := NewServerWithRouter(mockServer, router)
	mockServer.EXPECT().ListenAndServe().Return(nil)
	err := s.Start(context.Background())
	require.NoError(t, err)
}

func TestServerWithRouter_Start_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServer := NewMockHTTPServer(ctrl)
	router := chi.NewRouter()
	s := NewServerWithRouter(mockServer, router)
	expectedErr := errors.New("some error")
	mockServer.EXPECT().ListenAndServe().Return(expectedErr)
	err := s.Start(context.Background())
	assert.EqualError(t, err, expectedErr.Error())
}

func TestServerWithRouter_Start_ServerClosed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServer := NewMockHTTPServer(ctrl)
	router := chi.NewRouter()
	s := NewServerWithRouter(mockServer, router)
	mockServer.EXPECT().ListenAndServe().Return(http.ErrServerClosed)
	err := s.Start(context.Background())
	assert.NoError(t, err)
}

func TestServerWithRouter_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServer := NewMockHTTPServer(ctrl)
	router := chi.NewRouter()
	s := NewServerWithRouter(mockServer, router)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	mockServer.EXPECT().Shutdown(gomock.Any()).Return(nil)
	err := s.Stop(ctx)
	assert.NoError(t, err)
}

func TestServerWithRouter_Stop_WithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServer := NewMockHTTPServer(ctrl)
	router := chi.NewRouter()
	s := NewServerWithRouter(mockServer, router)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	mockServer.EXPECT().Shutdown(gomock.Any()).Return(errors.New("shutdown error"))
	err := s.Stop(ctx)
	assert.EqualError(t, err, "shutdown error")
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
}
