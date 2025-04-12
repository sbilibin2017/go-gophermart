package usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/usecases"
	"github.com/stretchr/testify/assert"
)

func TestUserRegisterUsecase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLV := usecases.NewMockLoginValidator(ctrl)
	mockPV := usecases.NewMockPasswordValidator(ctrl)
	mockService := usecases.NewMockUserRegisterService(ctrl)

	req := &usecases.UserRegisterRequest{
		Login:    "testuser",
		Password: "testpass123",
	}

	expectedToken := "mocked-jwt-token"

	mockLV.EXPECT().Validate(req.Login).Return(nil)
	mockPV.EXPECT().Validate(req.Password).Return(nil)
	mockService.EXPECT().Register(gomock.Any(), &services.User{
		Login:    req.Login,
		Password: req.Password,
	}).Return(expectedToken, nil)

	usecase := usecases.NewUserRegisterUsecase(mockLV, mockPV, mockService)
	resp, err := usecase.Execute(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedToken, resp.AccessToken)
}

func TestUserRegisterUsecase_Execute_LoginValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLV := usecases.NewMockLoginValidator(ctrl)
	mockPV := usecases.NewMockPasswordValidator(ctrl)
	mockService := usecases.NewMockUserRegisterService(ctrl)

	req := &usecases.UserRegisterRequest{
		Login:    "bad-login",
		Password: "password",
	}

	expectedErr := errors.New("invalid login")

	mockLV.EXPECT().Validate(req.Login).Return(expectedErr)

	usecase := usecases.NewUserRegisterUsecase(mockLV, mockPV, mockService)
	resp, err := usecase.Execute(context.Background(), req)

	assert.Nil(t, resp)
	assert.EqualError(t, err, expectedErr.Error())
}

func TestUserRegisterUsecase_Execute_PasswordValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLV := usecases.NewMockLoginValidator(ctrl)
	mockPV := usecases.NewMockPasswordValidator(ctrl)
	mockService := usecases.NewMockUserRegisterService(ctrl)

	req := &usecases.UserRegisterRequest{
		Login:    "user",
		Password: "bad",
	}

	expectedErr := errors.New("weak password")

	mockLV.EXPECT().Validate(req.Login).Return(nil)
	mockPV.EXPECT().Validate(req.Password).Return(expectedErr)

	usecase := usecases.NewUserRegisterUsecase(mockLV, mockPV, mockService)
	resp, err := usecase.Execute(context.Background(), req)

	assert.Nil(t, resp)
	assert.EqualError(t, err, expectedErr.Error())
}

func TestUserRegisterUsecase_Execute_RegisterError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLV := usecases.NewMockLoginValidator(ctrl)
	mockPV := usecases.NewMockPasswordValidator(ctrl)
	mockService := usecases.NewMockUserRegisterService(ctrl)

	req := &usecases.UserRegisterRequest{
		Login:    "user",
		Password: "strongpass",
	}

	expectedErr := errors.New("registration failed")

	mockLV.EXPECT().Validate(req.Login).Return(nil)
	mockPV.EXPECT().Validate(req.Password).Return(nil)
	mockService.EXPECT().Register(gomock.Any(), &services.User{
		Login:    req.Login,
		Password: req.Password,
	}).Return("", expectedErr)

	usecase := usecases.NewUserRegisterUsecase(mockLV, mockPV, mockService)
	resp, err := usecase.Execute(context.Background(), req)

	assert.Nil(t, resp)
	assert.EqualError(t, err, expectedErr.Error())
}
