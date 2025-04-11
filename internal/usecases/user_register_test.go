package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/services"

	"github.com/stretchr/testify/assert"
)

func TestUserRegisterUsecase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки для зависимостей
	mockLoginValidator := NewMockLoginValidator(ctrl)
	mockPasswordValidator := NewMockPasswordValidator(ctrl)
	mockUserRegisterService := NewMockUserRegisterService(ctrl)

	// Настроим моки
	mockLoginValidator.EXPECT().Validate("validLogin").Return(nil).Times(1)
	mockPasswordValidator.EXPECT().Validate("validPassword").Return(nil).Times(1)
	mockUserRegisterService.EXPECT().Register(context.Background(), &services.User{Login: "validLogin", Password: "validPassword"}).
		Return("accessToken", nil).Times(1)

	// Создаем экземпляр UserRegisterUsecase
	uc := &UserRegisterUsecase{
		lv:  mockLoginValidator,
		pv:  mockPasswordValidator,
		svc: mockUserRegisterService,
	}

	// Создаем запрос
	req := &UserRegisterRequest{
		Login:    "validLogin",
		Password: "validPassword",
	}

	// Выполняем метод Execute
	resp, err := uc.Execute(context.Background(), req)

	// Проверяем, что ошибок нет и токен возвращен
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "accessToken", resp.AccessToken)
}

func TestUserRegisterUsecase_Execute_LoginValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки для зависимостей
	mockLoginValidator := NewMockLoginValidator(ctrl)
	mockPasswordValidator := NewMockPasswordValidator(ctrl)
	mockUserRegisterService := NewMockUserRegisterService(ctrl)

	// Настроим моки
	mockLoginValidator.EXPECT().Validate("invalidLogin").Return(errors.New("invalid login")).Times(1)

	// Создаем экземпляр UserRegisterUsecase
	uc := &UserRegisterUsecase{
		lv:  mockLoginValidator,
		pv:  mockPasswordValidator,
		svc: mockUserRegisterService,
	}

	// Создаем запрос с неправильным логином
	req := &UserRegisterRequest{
		Login:    "invalidLogin",
		Password: "validPassword",
	}

	// Выполняем метод Execute
	resp, err := uc.Execute(context.Background(), req)

	// Проверяем, что ошибка вернулась, и ответ пустой
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "invalid login")
}

func TestUserRegisterUsecase_Execute_PasswordValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки для зависимостей
	mockLoginValidator := NewMockLoginValidator(ctrl)
	mockPasswordValidator := NewMockPasswordValidator(ctrl)
	mockUserRegisterService := NewMockUserRegisterService(ctrl)

	// Настроим моки
	mockLoginValidator.EXPECT().Validate("validLogin").Return(nil).Times(1)
	mockPasswordValidator.EXPECT().Validate("invalidPassword").Return(errors.New("invalid password")).Times(1)

	// Создаем экземпляр UserRegisterUsecase
	uc := &UserRegisterUsecase{
		lv:  mockLoginValidator,
		pv:  mockPasswordValidator,
		svc: mockUserRegisterService,
	}

	// Создаем запрос с неправильным паролем
	req := &UserRegisterRequest{
		Login:    "validLogin",
		Password: "invalidPassword",
	}

	// Выполняем метод Execute
	resp, err := uc.Execute(context.Background(), req)

	// Проверяем, что ошибка вернулась, и ответ пустой
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "invalid password")
}

func TestUserRegisterUsecase_Execute_UserRegisterServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки для зависимостей
	mockLoginValidator := NewMockLoginValidator(ctrl)
	mockPasswordValidator := NewMockPasswordValidator(ctrl)
	mockUserRegisterService := NewMockUserRegisterService(ctrl)

	// Настроим моки
	mockLoginValidator.EXPECT().Validate("validLogin").Return(nil).Times(1)
	mockPasswordValidator.EXPECT().Validate("validPassword").Return(nil).Times(1)
	mockUserRegisterService.EXPECT().Register(context.Background(), &services.User{Login: "validLogin", Password: "validPassword"}).
		Return("", errors.New("service error")).Times(1)

	// Создаем экземпляр UserRegisterUsecase
	uc := &UserRegisterUsecase{
		lv:  mockLoginValidator,
		pv:  mockPasswordValidator,
		svc: mockUserRegisterService,
	}

	// Создаем запрос
	req := &UserRegisterRequest{
		Login:    "validLogin",
		Password: "validPassword",
	}

	// Выполняем метод Execute
	resp, err := uc.Execute(context.Background(), req)

	// Проверяем, что ошибка вернулась и ответ пустой
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "service error")
}
