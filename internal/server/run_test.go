package server

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRunServerStartedSuccessfully(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServer := NewMockServer(ctrl)
	mockServer.EXPECT().ListenAndServe().Return(nil).Times(1)
	mockServer.EXPECT().Shutdown(gomock.Any()).Times(1)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, mockServer)
	}()
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out")
	}
}

func TestRunServerFailedToStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServer := NewMockServer(ctrl)
	expectedErr := fmt.Errorf("ListenAndServe error")
	mockServer.EXPECT().ListenAndServe().Return(expectedErr).Times(1)
	mockServer.EXPECT().Shutdown(gomock.Any()).Times(0)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, mockServer)
	}()
	select {
	case err := <-done:
		require.EqualError(t, err, expectedErr.Error())
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out")
	}
}

func TestRunServerGracefulShutdown(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServer := NewMockServer(ctrl)
	mockServer.EXPECT().ListenAndServe().Return(nil).Times(1)
	mockServer.EXPECT().Shutdown(gomock.Any()).Return(nil).Times(1)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, mockServer)
	}()
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out")
	}
}

func TestRunShutdownError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServer := NewMockServer(ctrl)
	mockServer.EXPECT().ListenAndServe().Return(nil).Times(1)
	expectedShutdownErr := fmt.Errorf("Shutdown error")
	mockServer.EXPECT().Shutdown(gomock.Any()).Return(expectedShutdownErr).Times(1)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, mockServer)
	}()
	select {
	case err := <-done:
		require.EqualError(t, err, expectedShutdownErr.Error())
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out")
	}
}
