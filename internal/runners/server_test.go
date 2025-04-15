package runners

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestServerShutdownCorrectly(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := NewMockServer(ctrl)

	ctx, cancel := context.WithCancel(context.Background())

	mockServer.EXPECT().ListenAndServe().DoAndReturn(func() error {
		<-ctx.Done()
		return nil
	}).Times(1)

	mockServer.EXPECT().Shutdown(gomock.Any()).Return(nil).Times(1)

	done := make(chan struct{})

	go func() {
		err := RunServer(ctx, mockServer)
		assert.NoError(t, err)
		close(done)
	}()

	time.Sleep(200 * time.Millisecond)
	cancel()
	<-done
}

func TestServerShutdownError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := NewMockServer(ctrl)

	ctx, cancel := context.WithCancel(context.Background())

	shutdownErr := errors.New("shutdown error")

	mockServer.EXPECT().ListenAndServe().DoAndReturn(func() error {
		<-ctx.Done()
		return nil
	}).Times(1)

	mockServer.EXPECT().Shutdown(gomock.Any()).Return(shutdownErr).Times(1)

	done := make(chan struct{})

	go func() {
		err := RunServer(ctx, mockServer)
		assert.Error(t, err)
		assert.Equal(t, shutdownErr, err)
		close(done)
	}()

	time.Sleep(200 * time.Millisecond)
	cancel()
	<-done
}

func TestRunServer_ReturnsListenAndServeErrorImmediately(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := NewMockServer(ctrl)

	expectedErr := errors.New("unexpected shutdown")

	mockServer.EXPECT().ListenAndServe().Return(expectedErr).Times(1)

	ctx := context.Background()

	err := RunServer(ctx, mockServer)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}
