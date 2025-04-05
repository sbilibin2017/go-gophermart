// Code generated by MockGen. DO NOT EDIT.
// Source: /home/sergey/Go/go-gophermart/go-gophermart/cmd/accrual/app/container.go

// Package app is a generated GoMock package.
package app

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockContainerConfig is a mock of ContainerConfig interface.
type MockContainerConfig struct {
	ctrl     *gomock.Controller
	recorder *MockContainerConfigMockRecorder
}

// MockContainerConfigMockRecorder is the mock recorder for MockContainerConfig.
type MockContainerConfigMockRecorder struct {
	mock *MockContainerConfig
}

// NewMockContainerConfig creates a new mock instance.
func NewMockContainerConfig(ctrl *gomock.Controller) *MockContainerConfig {
	mock := &MockContainerConfig{ctrl: ctrl}
	mock.recorder = &MockContainerConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContainerConfig) EXPECT() *MockContainerConfigMockRecorder {
	return m.recorder
}

// GetDatabaseURI mocks base method.
func (m *MockContainerConfig) GetDatabaseURI() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDatabaseURI")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetDatabaseURI indicates an expected call of GetDatabaseURI.
func (mr *MockContainerConfigMockRecorder) GetDatabaseURI() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDatabaseURI", reflect.TypeOf((*MockContainerConfig)(nil).GetDatabaseURI))
}

// GetRunAddress mocks base method.
func (m *MockContainerConfig) GetRunAddress() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRunAddress")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetRunAddress indicates an expected call of GetRunAddress.
func (mr *MockContainerConfigMockRecorder) GetRunAddress() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRunAddress", reflect.TypeOf((*MockContainerConfig)(nil).GetRunAddress))
}
