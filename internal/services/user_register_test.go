package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"

	"github.com/stretchr/testify/assert"
)

func TestUserRegisterService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserFilterRepo := NewMockUserFilterRepo(ctrl)
	mockUserSaveRepo := NewMockUserSaveRepo(ctrl)
	mockHasher := NewMockHasher(ctrl)
	mockJWTGenerator := NewMockJWTGenerator(ctrl)

	service := NewUserRegisterService(mockUserFilterRepo, mockUserSaveRepo, mockHasher, mockJWTGenerator)

	t.Run("success", func(t *testing.T) {
		// Arrange
		user := &User{Login: "testuser", Password: "password123"}
		mockUserFilterRepo.EXPECT().Filter(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
		mockHasher.EXPECT().Hash(user.Password).Return("hashedpassword").Times(1)
		mockUserSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		mockJWTGenerator.EXPECT().Generate(user.Login).Return("mocked-jwt-token").Times(1)

		token, err := service.Register(context.Background(), user)

		assert.NoError(t, err)
		assert.Equal(t, "mocked-jwt-token", token)
	})

	t.Run("user already exists", func(t *testing.T) {
		user := &User{Login: "existinguser", Password: "password123"}
		mockUserFilterRepo.EXPECT().Filter(gomock.Any(), gomock.Any()).Return(&repositories.UserFiltered{}, nil).Times(1)

		token, err := service.Register(context.Background(), user)

		assert.Error(t, err)
		assert.Equal(t, err, ErrUserAlreadyExists)
		assert.Empty(t, token)
	})

	t.Run("error filtering user", func(t *testing.T) {
		user := &User{Login: "testuser", Password: "password123"}
		mockUserFilterRepo.EXPECT().Filter(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error")).Times(1)

		token, err := service.Register(context.Background(), user)

		assert.Error(t, err)
		assert.Equal(t, err, ErrUserRegisterInternal)
		assert.Empty(t, token)
	})

	t.Run("error saving user", func(t *testing.T) {
		user := &User{Login: "testuser", Password: "password123"}
		mockUserFilterRepo.EXPECT().Filter(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
		mockHasher.EXPECT().Hash(user.Password).Return("hashedpassword").Times(1)
		mockUserSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errors.New("save error")).Times(1)

		token, err := service.Register(context.Background(), user)

		assert.Error(t, err)
		assert.Equal(t, err, ErrUserRegisterInternal)
		assert.Empty(t, token)
	})
}
