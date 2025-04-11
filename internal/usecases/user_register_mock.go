// Code generated by MockGen. DO NOT EDIT.
// Source: /home/sergey/Go/go-gophermart/internal/usecases/user_register.go

// Package usecases is a generated GoMock package.
package usecases

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/sbilibin2017/go-gophermart/internal/domain"
)

// MockLoginValidator is a mock of LoginValidator interface.
type MockLoginValidator struct {
	ctrl     *gomock.Controller
	recorder *MockLoginValidatorMockRecorder
}

// MockLoginValidatorMockRecorder is the mock recorder for MockLoginValidator.
type MockLoginValidatorMockRecorder struct {
	mock *MockLoginValidator
}

// NewMockLoginValidator creates a new mock instance.
func NewMockLoginValidator(ctrl *gomock.Controller) *MockLoginValidator {
	mock := &MockLoginValidator{ctrl: ctrl}
	mock.recorder = &MockLoginValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLoginValidator) EXPECT() *MockLoginValidatorMockRecorder {
	return m.recorder
}

// Validate mocks base method.
func (m *MockLoginValidator) Validate(login string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", login)
	ret0, _ := ret[0].(error)
	return ret0
}

// Validate indicates an expected call of Validate.
func (mr *MockLoginValidatorMockRecorder) Validate(login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockLoginValidator)(nil).Validate), login)
}

// MockPasswordValidator is a mock of PasswordValidator interface.
type MockPasswordValidator struct {
	ctrl     *gomock.Controller
	recorder *MockPasswordValidatorMockRecorder
}

// MockPasswordValidatorMockRecorder is the mock recorder for MockPasswordValidator.
type MockPasswordValidatorMockRecorder struct {
	mock *MockPasswordValidator
}

// NewMockPasswordValidator creates a new mock instance.
func NewMockPasswordValidator(ctrl *gomock.Controller) *MockPasswordValidator {
	mock := &MockPasswordValidator{ctrl: ctrl}
	mock.recorder = &MockPasswordValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPasswordValidator) EXPECT() *MockPasswordValidatorMockRecorder {
	return m.recorder
}

// Validate mocks base method.
func (m *MockPasswordValidator) Validate(password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", password)
	ret0, _ := ret[0].(error)
	return ret0
}

// Validate indicates an expected call of Validate.
func (mr *MockPasswordValidatorMockRecorder) Validate(password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockPasswordValidator)(nil).Validate), password)
}

// MockUserRegisterService is a mock of UserRegisterService interface.
type MockUserRegisterService struct {
	ctrl     *gomock.Controller
	recorder *MockUserRegisterServiceMockRecorder
}

// MockUserRegisterServiceMockRecorder is the mock recorder for MockUserRegisterService.
type MockUserRegisterServiceMockRecorder struct {
	mock *MockUserRegisterService
}

// NewMockUserRegisterService creates a new mock instance.
func NewMockUserRegisterService(ctrl *gomock.Controller) *MockUserRegisterService {
	mock := &MockUserRegisterService{ctrl: ctrl}
	mock.recorder = &MockUserRegisterServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRegisterService) EXPECT() *MockUserRegisterServiceMockRecorder {
	return m.recorder
}

// Register mocks base method.
func (m *MockUserRegisterService) Register(ctx context.Context, u *domain.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, u)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockUserRegisterServiceMockRecorder) Register(ctx, u interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockUserRegisterService)(nil).Register), ctx, u)
}
