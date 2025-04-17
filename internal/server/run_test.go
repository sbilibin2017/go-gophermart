package server_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/log"
	"github.com/sbilibin2017/go-gophermart/internal/server"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.Init()
}

func TestRun(t *testing.T) {
	tests := []struct {
		name              string
		isErr             bool
		mockExpectActions func(ctrl *gomock.Controller, mockServer *server.MockServer)
	}{
		{
			name:  "TestRunSuccess",
			isErr: false,
			mockExpectActions: func(ctrl *gomock.Controller, mockServer *server.MockServer) {
				mockServer.EXPECT().ListenAndServe().Return(nil).Times(1)
				mockServer.EXPECT().Shutdown(gomock.Any()).Return(nil).Times(1)
			},
		},
		{
			name:  "TestRunListenAndServeError",
			isErr: true,
			mockExpectActions: func(ctrl *gomock.Controller, mockServer *server.MockServer) {
				listenAndServeError := errors.New("ListenAndServe error")
				mockServer.EXPECT().ListenAndServe().Return(listenAndServeError).Times(1)
			},
		},
		{
			name:  "TestRunShutdownError",
			isErr: true,
			mockExpectActions: func(ctrl *gomock.Controller, mockServer *server.MockServer) {
				mockServer.EXPECT().ListenAndServe().Return(nil).Times(1)
				shutdownError := errors.New("Shutdown error")
				mockServer.EXPECT().Shutdown(gomock.Any()).Return(shutdownError).Times(1)
			},
		},
		{
			name:  "TestRunContextCanceled",
			isErr: false,
			mockExpectActions: func(ctrl *gomock.Controller, mockServer *server.MockServer) {
				mockServer.EXPECT().ListenAndServe().Return(nil).Times(1)
				mockServer.EXPECT().Shutdown(gomock.Any()).Return(nil).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Mock the Server
			mockServer := server.NewMockServer(ctrl)

			// Set up the mock behavior for this test
			tt.mockExpectActions(ctrl, mockServer)

			// Create a context with a timeout to prevent the test from hanging indefinitely
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			// Run the function being tested
			err := server.Run(ctx, mockServer)

			// Assert that the correct behavior occurs
			if tt.isErr {
				assert.Error(t, err, "Expected an error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
			}
		})
	}
}
