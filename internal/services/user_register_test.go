package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestUserRegisterService_Register(t *testing.T) {
	type args struct {
		req *types.UserRegisterRequest
	}

	tests := []struct {
		name       string
		setupMocks func(
			v *services.MockUserRegisterValidator,
			uf *services.MockUserRegisterUserFilterOneRepository,
			us *services.MockUserRegisterUserSaveRepository,
		)
		args        args
		wantErr     *types.APIStatus
		wantSuccess *types.APIStatus
		expectToken bool
	}{
		{
			name: "success",
			args: args{
				req: &types.UserRegisterRequest{
					Login:    "newuser",
					Password: "password123",
				},
			},
			setupMocks: func(v *services.MockUserRegisterValidator, uf *services.MockUserRegisterUserFilterOneRepository, us *services.MockUserRegisterUserSaveRepository) {
				v.EXPECT().Struct(gomock.Any()).Return(nil)
				uf.EXPECT().FilterOne(gomock.Any(), "newuser").Return(nil, nil)
				us.EXPECT().Save(gomock.Any(), "newuser", gomock.Any()).Return(nil)
			},
			wantErr:     nil,
			wantSuccess: &types.APIStatus{StatusCode: 201, Message: "User successfully registered"},
			expectToken: true,
		},
		{
			name: "error on FilterOne",
			args: args{
				req: &types.UserRegisterRequest{
					Login:    "dbfail",
					Password: "password",
				},
			},
			setupMocks: func(v *services.MockUserRegisterValidator, uf *services.MockUserRegisterUserFilterOneRepository, us *services.MockUserRegisterUserSaveRepository) {
				v.EXPECT().Struct(gomock.Any()).Return(nil)
				uf.EXPECT().FilterOne(gomock.Any(), "dbfail").Return(nil, errors.New("db error"))
			},
			wantErr:     &types.APIStatus{StatusCode: 500, Message: "Internal server error while registering user"},
			wantSuccess: nil,
			expectToken: false,
		},
		{
			name: "validation error",
			args: args{
				req: &types.UserRegisterRequest{
					Login:    "",
					Password: "123",
				},
			},
			setupMocks: func(v *services.MockUserRegisterValidator, uf *services.MockUserRegisterUserFilterOneRepository, us *services.MockUserRegisterUserSaveRepository) {
				v.EXPECT().Struct(gomock.Any()).Return(errors.New("validation failed"))
			},
			wantErr:     &types.APIStatus{StatusCode: 400, Message: "Invalid request"},
			wantSuccess: nil,
			expectToken: false,
		},
		{
			name: "user already exists",
			args: args{
				req: &types.UserRegisterRequest{
					Login:    "existing",
					Password: "password123",
				},
			},
			setupMocks: func(v *services.MockUserRegisterValidator, uf *services.MockUserRegisterUserFilterOneRepository, us *services.MockUserRegisterUserSaveRepository) {
				v.EXPECT().Struct(gomock.Any()).Return(nil)
				uf.EXPECT().FilterOne(gomock.Any(), "existing").Return(&types.UserDB{Login: "existing"}, nil)
			},
			wantErr:     &types.APIStatus{StatusCode: 409, Message: "User login already exists"},
			wantSuccess: nil,
			expectToken: false,
		},
		{
			name: "internal error on Save",
			args: args{
				req: &types.UserRegisterRequest{
					Login:    "fail",
					Password: "pass",
				},
			},
			setupMocks: func(v *services.MockUserRegisterValidator, uf *services.MockUserRegisterUserFilterOneRepository, us *services.MockUserRegisterUserSaveRepository) {
				v.EXPECT().Struct(gomock.Any()).Return(nil)
				uf.EXPECT().FilterOne(gomock.Any(), "fail").Return(nil, nil)
				us.EXPECT().Save(gomock.Any(), "fail", gomock.Any()).Return(errors.New("db error"))
			},
			wantErr:     &types.APIStatus{StatusCode: 500, Message: "Internal server error while registering user"},
			wantSuccess: nil,
			expectToken: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockValidator := services.NewMockUserRegisterValidator(ctrl)
			mockFilter := services.NewMockUserRegisterUserFilterOneRepository(ctrl)
			mockSaver := services.NewMockUserRegisterUserSaveRepository(ctrl)

			tt.setupMocks(mockValidator, mockFilter, mockSaver)

			svc := services.NewUserRegisterService(
				mockValidator,
				mockFilter,
				mockSaver,
				"secret",
				time.Hour,
				"gophermart",
			)

			token, success, err := svc.Register(context.Background(), tt.args.req)

			assert.Equal(t, tt.wantSuccess, success)
			assert.Equal(t, tt.wantErr, err)

			if tt.expectToken {
				assert.NotEmpty(t, token)
			} else {
				assert.Empty(t, token)
			}
		})
	}
}
