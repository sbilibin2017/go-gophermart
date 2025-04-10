package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestUserRegisterService_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock instances
	mockUserGetRepo := NewMockUserGetRepo(ctrl)
	mockUserSaveRepo := NewMockUserSaveRepo(ctrl)
	mockUnitOfWork := NewMockUnitOfWork(ctrl)

	// Create the service
	service := NewUserRegisterService(
		&configs.GophermartConfig{
			JWTSecretKey: "secret",
			JWTExp:       3600,
		},
		mockUserGetRepo,
		mockUserSaveRepo,
		mockUnitOfWork,
	)

	// Test data
	user := &domain.User{
		Login:    "testuser",
		Password: "password123",
	}

	// Define expected behavior for mock objects
	mockUserGetRepo.EXPECT().GetByParam(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
	mockUserSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	// Mock the UnitOfWork's Do method to invoke the passed function
	mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, operation func(tx *sql.Tx) error) error {
		// Call the provided operation (the function passed to Do)
		return operation(nil) // nil simulates the tx (transaction), but in real scenarios it would be an actual DB transaction object.
	}).Times(1)

	// Call Register
	token, err := service.Register(context.Background(), user)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestUserRegisterService_Register_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock instances
	mockUserGetRepo := NewMockUserGetRepo(ctrl)
	mockUserSaveRepo := NewMockUserSaveRepo(ctrl)
	mockUnitOfWork := NewMockUnitOfWork(ctrl)

	// Create the service
	service := NewUserRegisterService(
		&configs.GophermartConfig{
			JWTSecretKey: "secret",
			JWTExp:       3600,
		},
		mockUserGetRepo,
		mockUserSaveRepo,
		mockUnitOfWork,
	)

	// Test data
	user := &domain.User{
		Login:    "existinguser",
		Password: "password123",
	}

	// Define expected behavior for mock objects
	mockUserGetRepo.EXPECT().GetByParam(gomock.Any(), gomock.Any()).Return(map[string]any{"login": "existinguser"}, nil).Times(1)

	// Mock the UnitOfWork's Do method
	mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, operation func(tx *sql.Tx) error) error {
		// Call the provided operation (the function passed to Do)
		return operation(nil) // Simulate transaction context
	}).Times(1)

	// Call Register
	token, err := service.Register(context.Background(), user)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, ErrUserAlreadyExists, err)
	assert.Empty(t, token)
}

func TestUserRegisterService_Register_InternalErrorOnSave(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock instances
	mockUserGetRepo := NewMockUserGetRepo(ctrl)
	mockUserSaveRepo := NewMockUserSaveRepo(ctrl)
	mockUnitOfWork := NewMockUnitOfWork(ctrl)

	// Create the service
	service := NewUserRegisterService(
		&configs.GophermartConfig{
			JWTSecretKey: "secret",
			JWTExp:       3600,
		},
		mockUserGetRepo,
		mockUserSaveRepo,
		mockUnitOfWork,
	)

	// Test data
	user := &domain.User{
		Login:    "newuser",
		Password: "password123",
	}

	// Define expected behavior for mock objects
	mockUserGetRepo.EXPECT().GetByParam(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
	mockUserSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errors.New("internal error")).Times(1)

	// Mock the UnitOfWork's Do method
	mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, operation func(tx *sql.Tx) error) error {
		// Call the provided operation (the function passed to Do)
		return operation(nil) // Simulate transaction context
	}).Times(1)

	// Call Register
	token, err := service.Register(context.Background(), user)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, ErrUserRegisterInternal, err)
	assert.Empty(t, token)
}

func TestUserRegisterService_GetByParam_ErrorHandling(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock instances
	mockUserGetRepo := NewMockUserGetRepo(ctrl)
	mockUserSaveRepo := NewMockUserSaveRepo(ctrl)
	mockUnitOfWork := NewMockUnitOfWork(ctrl)

	// Create the service
	service := NewUserRegisterService(
		&configs.GophermartConfig{
			JWTSecretKey: "secret",
			JWTExp:       3600,
		},
		mockUserGetRepo,
		mockUserSaveRepo,
		mockUnitOfWork,
	)

	// Test data
	user := &domain.User{
		Login:    "newuser",
		Password: "password123",
	}

	// Mock the behavior of GetByParam to return an error
	mockUserGetRepo.EXPECT().GetByParam(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error")).Times(1)

	// Mock the behavior of UnitOfWork's Do method to call GetByParam and return nil
	mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, operation func(tx *sql.Tx) error) error {
		// Simulate the behavior of GetByParam being called inside the operation
		return operation(nil) // Passing nil as tx since we're not focusing on database transactions here
	}).Times(1)

	// Call Register method
	token, err := service.Register(context.Background(), user)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, ErrUserRegisterInternal, err)
	assert.Empty(t, token)
}
