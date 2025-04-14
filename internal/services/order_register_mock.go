// Code generated by MockGen. DO NOT EDIT.
// Source: /home/sergey/Go/go-gophermart/internal/services/order_register.go

// Package services is a generated GoMock package.
package services

import (
	context "context"
	sql "database/sql"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	types "github.com/sbilibin2017/go-gophermart/internal/types"
)

// MockOrderRegisterOrderExistsRepository is a mock of OrderRegisterOrderExistsRepository interface.
type MockOrderRegisterOrderExistsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockOrderRegisterOrderExistsRepositoryMockRecorder
}

// MockOrderRegisterOrderExistsRepositoryMockRecorder is the mock recorder for MockOrderRegisterOrderExistsRepository.
type MockOrderRegisterOrderExistsRepositoryMockRecorder struct {
	mock *MockOrderRegisterOrderExistsRepository
}

// NewMockOrderRegisterOrderExistsRepository creates a new mock instance.
func NewMockOrderRegisterOrderExistsRepository(ctrl *gomock.Controller) *MockOrderRegisterOrderExistsRepository {
	mock := &MockOrderRegisterOrderExistsRepository{ctrl: ctrl}
	mock.recorder = &MockOrderRegisterOrderExistsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderRegisterOrderExistsRepository) EXPECT() *MockOrderRegisterOrderExistsRepositoryMockRecorder {
	return m.recorder
}

// Exists mocks base method.
func (m *MockOrderRegisterOrderExistsRepository) Exists(ctx context.Context, orderID *types.OrderExistsFilter) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists", ctx, orderID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists.
func (mr *MockOrderRegisterOrderExistsRepositoryMockRecorder) Exists(ctx, orderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockOrderRegisterOrderExistsRepository)(nil).Exists), ctx, orderID)
}

// MockOrderRegisterOrderSaveRepository is a mock of OrderRegisterOrderSaveRepository interface.
type MockOrderRegisterOrderSaveRepository struct {
	ctrl     *gomock.Controller
	recorder *MockOrderRegisterOrderSaveRepositoryMockRecorder
}

// MockOrderRegisterOrderSaveRepositoryMockRecorder is the mock recorder for MockOrderRegisterOrderSaveRepository.
type MockOrderRegisterOrderSaveRepositoryMockRecorder struct {
	mock *MockOrderRegisterOrderSaveRepository
}

// NewMockOrderRegisterOrderSaveRepository creates a new mock instance.
func NewMockOrderRegisterOrderSaveRepository(ctrl *gomock.Controller) *MockOrderRegisterOrderSaveRepository {
	mock := &MockOrderRegisterOrderSaveRepository{ctrl: ctrl}
	mock.recorder = &MockOrderRegisterOrderSaveRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderRegisterOrderSaveRepository) EXPECT() *MockOrderRegisterOrderSaveRepositoryMockRecorder {
	return m.recorder
}

// Save mocks base method.
func (m *MockOrderRegisterOrderSaveRepository) Save(ctx context.Context, order *types.OrderDB) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockOrderRegisterOrderSaveRepositoryMockRecorder) Save(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockOrderRegisterOrderSaveRepository)(nil).Save), ctx, order)
}

// MockOrderRegisterRewardFilterRepository is a mock of OrderRegisterRewardFilterRepository interface.
type MockOrderRegisterRewardFilterRepository struct {
	ctrl     *gomock.Controller
	recorder *MockOrderRegisterRewardFilterRepositoryMockRecorder
}

// MockOrderRegisterRewardFilterRepositoryMockRecorder is the mock recorder for MockOrderRegisterRewardFilterRepository.
type MockOrderRegisterRewardFilterRepositoryMockRecorder struct {
	mock *MockOrderRegisterRewardFilterRepository
}

// NewMockOrderRegisterRewardFilterRepository creates a new mock instance.
func NewMockOrderRegisterRewardFilterRepository(ctrl *gomock.Controller) *MockOrderRegisterRewardFilterRepository {
	mock := &MockOrderRegisterRewardFilterRepository{ctrl: ctrl}
	mock.recorder = &MockOrderRegisterRewardFilterRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderRegisterRewardFilterRepository) EXPECT() *MockOrderRegisterRewardFilterRepositoryMockRecorder {
	return m.recorder
}

// Filter mocks base method.
func (m *MockOrderRegisterRewardFilterRepository) Filter(ctx context.Context, filter []*types.RewardFilter) ([]*types.RewardDB, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Filter", ctx, filter)
	ret0, _ := ret[0].([]*types.RewardDB)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Filter indicates an expected call of Filter.
func (mr *MockOrderRegisterRewardFilterRepositoryMockRecorder) Filter(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Filter", reflect.TypeOf((*MockOrderRegisterRewardFilterRepository)(nil).Filter), ctx, filter)
}

// MockTx is a mock of Tx interface.
type MockTx struct {
	ctrl     *gomock.Controller
	recorder *MockTxMockRecorder
}

// MockTxMockRecorder is the mock recorder for MockTx.
type MockTxMockRecorder struct {
	mock *MockTx
}

// NewMockTx creates a new mock instance.
func NewMockTx(ctrl *gomock.Controller) *MockTx {
	mock := &MockTx{ctrl: ctrl}
	mock.recorder = &MockTxMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTx) EXPECT() *MockTxMockRecorder {
	return m.recorder
}

// Do mocks base method.
func (m *MockTx) Do(ctx context.Context, operation func(*sql.Tx) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", ctx, operation)
	ret0, _ := ret[0].(error)
	return ret0
}

// Do indicates an expected call of Do.
func (mr *MockTxMockRecorder) Do(ctx, operation interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockTx)(nil).Do), ctx, operation)
}
