// Code generated by MockGen. DO NOT EDIT.
// Source: /home/sergey/Go/go-gophermart/internal/repositories/reward_save.go

// Package repositories is a generated GoMock package.
package repositories

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRewardSaveExecutor is a mock of RewardSaveExecutor interface.
type MockRewardSaveExecutor struct {
	ctrl     *gomock.Controller
	recorder *MockRewardSaveExecutorMockRecorder
}

// MockRewardSaveExecutorMockRecorder is the mock recorder for MockRewardSaveExecutor.
type MockRewardSaveExecutorMockRecorder struct {
	mock *MockRewardSaveExecutor
}

// NewMockRewardSaveExecutor creates a new mock instance.
func NewMockRewardSaveExecutor(ctrl *gomock.Controller) *MockRewardSaveExecutor {
	mock := &MockRewardSaveExecutor{ctrl: ctrl}
	mock.recorder = &MockRewardSaveExecutorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRewardSaveExecutor) EXPECT() *MockRewardSaveExecutorMockRecorder {
	return m.recorder
}

// Execute mocks base method.
func (m *MockRewardSaveExecutor) Execute(ctx context.Context, query string, args map[string]any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", ctx, query, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// Execute indicates an expected call of Execute.
func (mr *MockRewardSaveExecutorMockRecorder) Execute(ctx, query, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockRewardSaveExecutor)(nil).Execute), ctx, query, args)
}
