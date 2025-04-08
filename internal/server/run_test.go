package server

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRun_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServer := NewMockServer(ctrl)
	mockServer.EXPECT().ListenAndServe().DoAndReturn(func() error {
		<-time.After(100 * time.Millisecond)
		return nil
	}).Times(1)
	mockServer.EXPECT().Shutdown(gomock.Any()).Return(nil).Times(1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := Run(ctx, mockServer)
		assert.NoError(t, err)
	}()
	time.Sleep(50 * time.Millisecond)
	cancel()
	wg.Wait()
}

func TestRun_ServerFailsToStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServer := NewMockServer(ctrl)
	mockServer.EXPECT().ListenAndServe().Return(errors.New("server error")).Times(1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := Run(ctx, mockServer)
		assert.EqualError(t, err, "server error")
	}()
	time.Sleep(50 * time.Millisecond)
	cancel()
	wg.Wait()
}

func TestRun_ServerShutdownWithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServer := NewMockServer(ctrl)
	mockServer.EXPECT().ListenAndServe().DoAndReturn(func() error {
		<-time.After(100 * time.Millisecond)
		return nil
	}).Times(1)
	mockServer.EXPECT().Shutdown(gomock.Any()).Return(errors.New("shutdown error")).Times(1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := Run(ctx, mockServer)
		assert.EqualError(t, err, "shutdown error")
	}()
	time.Sleep(50 * time.Millisecond)
	cancel()
	wg.Wait()
}
