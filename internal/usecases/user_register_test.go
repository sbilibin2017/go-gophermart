package usecases_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/usecases"
	"github.com/stretchr/testify/assert"
)

func TestUserRegisterUsecase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoginValidator := usecases.NewMockLoginValidator(ctrl)
	mockPasswordValidator := usecases.NewMockPasswordValidator(ctrl)
	mockUserRegisterService := usecases.NewMockUserRegisterService(ctrl)
	mockUnitOfWork := usecases.NewMockUnitOfWork(ctrl)

	login := "testuser"
	password := "testpassword"
	expectedToken := "mockedAccessToken"

	mockLoginValidator.EXPECT().Validate(login).Return(nil).Times(1)
	mockPasswordValidator.EXPECT().Validate(password).Return(nil).Times(1)
	mockUserRegisterService.EXPECT().Register(gomock.Any(), gomock.Any()).Return(expectedToken, nil).Times(1)
	mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, operation func(tx *sql.Tx) error) error {
		return operation(nil)
	}).Times(1)

	usecase := usecases.NewUserRegisterUsecase(mockUnitOfWork, mockLoginValidator, mockPasswordValidator, mockUserRegisterService)

	req := &usecases.UserRegisterRequest{Login: login, Password: password}
	resp, err := usecase.Execute(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedToken, resp.AccessToken)
}

func TestUserRegisterUsecase_Execute_LoginValidationFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoginValidator := usecases.NewMockLoginValidator(ctrl)
	mockPasswordValidator := usecases.NewMockPasswordValidator(ctrl)
	mockUserRegisterService := usecases.NewMockUserRegisterService(ctrl)
	mockUnitOfWork := usecases.NewMockUnitOfWork(ctrl)

	login := "testuser"
	password := "testpassword"

	mockLoginValidator.EXPECT().Validate(login).Return(errors.New("invalid login")).Times(1)
	mockPasswordValidator.EXPECT().Validate(password).Times(0)
	mockUserRegisterService.EXPECT().Register(gomock.Any(), gomock.Any()).Times(0)

	mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, operation func(tx *sql.Tx) error) error {
		return operation(nil)
	}).Times(1)

	usecase := usecases.NewUserRegisterUsecase(mockUnitOfWork, mockLoginValidator, mockPasswordValidator, mockUserRegisterService)

	req := &usecases.UserRegisterRequest{Login: login, Password: password}
	resp, err := usecase.Execute(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "invalid login", err.Error())
}

func TestUserRegisterUsecase_Execute_PasswordValidationFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoginValidator := usecases.NewMockLoginValidator(ctrl)
	mockPasswordValidator := usecases.NewMockPasswordValidator(ctrl)
	mockUserRegisterService := usecases.NewMockUserRegisterService(ctrl)
	mockUnitOfWork := usecases.NewMockUnitOfWork(ctrl)

	login := "testuser"
	password := "testpassword"

	mockLoginValidator.EXPECT().Validate(login).Return(nil).Times(1)
	mockPasswordValidator.EXPECT().Validate(password).Return(errors.New("invalid password")).Times(1)
	mockUserRegisterService.EXPECT().Register(gomock.Any(), gomock.Any()).Times(0)

	mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, operation func(tx *sql.Tx) error) error {
		return operation(nil)
	}).Times(1)

	usecase := usecases.NewUserRegisterUsecase(mockUnitOfWork, mockLoginValidator, mockPasswordValidator, mockUserRegisterService)

	req := &usecases.UserRegisterRequest{Login: login, Password: password}
	resp, err := usecase.Execute(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "invalid password", err.Error())
}

func TestUserRegisterUsecase_Execute_RegisterFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoginValidator := usecases.NewMockLoginValidator(ctrl)
	mockPasswordValidator := usecases.NewMockPasswordValidator(ctrl)
	mockUserRegisterService := usecases.NewMockUserRegisterService(ctrl)
	mockUnitOfWork := usecases.NewMockUnitOfWork(ctrl)

	login := "testuser"
	password := "testpassword"
	expectedError := errors.New("registration failed")

	mockLoginValidator.EXPECT().Validate(login).Return(nil).Times(1)
	mockPasswordValidator.EXPECT().Validate(password).Return(nil).Times(1)
	mockUserRegisterService.EXPECT().Register(gomock.Any(), gomock.Any()).Return("", expectedError).Times(1)

	mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, operation func(tx *sql.Tx) error) error {
		return operation(nil)
	}).Times(1)

	usecase := usecases.NewUserRegisterUsecase(mockUnitOfWork, mockLoginValidator, mockPasswordValidator, mockUserRegisterService)

	req := &usecases.UserRegisterRequest{Login: login, Password: password}
	resp, err := usecase.Execute(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "registration failed", err.Error())
}
