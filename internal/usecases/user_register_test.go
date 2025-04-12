package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/services"

	"github.com/stretchr/testify/assert"
)

func TestUserRegisterUsecase_Execute(t *testing.T) {
	type testCase struct {
		name       string
		login      string
		password   string
		setupMocks func(
			mockLoginValidator *MockLoginValidator,
			mockPasswordValidator *MockPasswordValidator,
			mockUserRegisterService *MockUserRegisterService,
		)
		expectedToken string
		expectedError error
	}

	testCases := []testCase{
		{
			name:     "Success",
			login:    "validLogin",
			password: "validPassword",
			setupMocks: func(mLogin *MockLoginValidator, mPass *MockPasswordValidator, mService *MockUserRegisterService) {
				mLogin.EXPECT().Validate("validLogin").Return(nil)
				mPass.EXPECT().Validate("validPassword").Return(nil)
				mService.EXPECT().
					Register(gomock.Any(), &services.User{Login: "validLogin", Password: "validPassword"}).
					Return("accessToken", nil)
			},
			expectedToken: "accessToken",
			expectedError: nil,
		},
		{
			name:     "Login validation error",
			login:    "invalidLogin",
			password: "validPassword",
			setupMocks: func(mLogin *MockLoginValidator, mPass *MockPasswordValidator, mService *MockUserRegisterService) {
				mLogin.EXPECT().Validate("invalidLogin").Return(errors.New("invalid login"))
			},
			expectedToken: "",
			expectedError: errors.New("invalid login"),
		},
		{
			name:     "Password validation error",
			login:    "validLogin",
			password: "invalidPassword",
			setupMocks: func(mLogin *MockLoginValidator, mPass *MockPasswordValidator, mService *MockUserRegisterService) {
				mLogin.EXPECT().Validate("validLogin").Return(nil)
				mPass.EXPECT().Validate("invalidPassword").Return(errors.New("invalid password"))
			},
			expectedToken: "",
			expectedError: errors.New("invalid password"),
		},
		{
			name:     "Register service error",
			login:    "validLogin",
			password: "validPassword",
			setupMocks: func(mLogin *MockLoginValidator, mPass *MockPasswordValidator, mService *MockUserRegisterService) {
				mLogin.EXPECT().Validate("validLogin").Return(nil)
				mPass.EXPECT().Validate("validPassword").Return(nil)
				mService.EXPECT().
					Register(gomock.Any(), &services.User{Login: "validLogin", Password: "validPassword"}).
					Return("", errors.New("service error"))
			},
			expectedToken: "",
			expectedError: errors.New("service error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLoginValidator := NewMockLoginValidator(ctrl)
			mockPasswordValidator := NewMockPasswordValidator(ctrl)
			mockUserRegisterService := NewMockUserRegisterService(ctrl)

			tc.setupMocks(mockLoginValidator, mockPasswordValidator, mockUserRegisterService)

			uc := NewUserRegisterUsecase(mockLoginValidator, mockPasswordValidator, mockUserRegisterService)

			req := &UserRegisterRequest{
				Login:    tc.login,
				Password: tc.password,
			}

			resp, err := uc.Execute(context.Background(), req)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tc.expectedToken, resp.AccessToken)
			}
		})
	}
}
