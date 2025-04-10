package services

import (
	"context"
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestUserService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserGetRepo := NewMockUserGetRepo(ctrl)
	mockUserSaveRepo := NewMockUserSaveRepo(ctrl)
	mockUnitOfWork := NewMockUnitOfWork(ctrl)
	mockHasher := NewMockHasher(ctrl)
	mockJWTEncoder := NewMockJWTEncoder(ctrl)

	svc := NewUserRegisterService(
		mockUserGetRepo,
		mockUserSaveRepo,
		mockUnitOfWork,
		mockHasher,
		mockJWTEncoder,
	)

	tests := []struct {
		name          string
		req           *domain.User
		mockSetup     func()
		expectedError error
		expectedResp  *domain.UserToken
	}{
		{
			name: "Happy path",
			req: &domain.User{
				Login:    "user1",
				Password: "Password123!",
			},
			mockSetup: func() {
				mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx *sql.Tx) error) error {
						return fn(nil)
					},
				)
				mockUserGetRepo.EXPECT().GetByParam(gomock.Any(), gomock.Any()).Return(nil, nil)
				mockHasher.EXPECT().Hash("Password123!").Return("hashedpassword123", nil)
				mockUserSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
				mockJWTEncoder.EXPECT().Encode("user1").Return("access-token", nil)
			},
			expectedError: nil,
			expectedResp: &domain.UserToken{
				Access: "access-token",
			},
		},
		{
			name: "User already exists",
			req: &domain.User{
				Login:    "user1",
				Password: "Password123!",
			},
			mockSetup: func() {
				mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx *sql.Tx) error) error {
						return fn(nil)
					},
				)
				mockUserGetRepo.EXPECT().GetByParam(gomock.Any(), gomock.Any()).Return(map[string]any{"login": "user1"}, nil)
				mockHasher.EXPECT().Hash(gomock.Any()).Times(0)
			},
			expectedError: ErrUserAlreadyExists,
			expectedResp:  nil,
		},
		{
			name: "Password hash failure",
			req: &domain.User{
				Login:    "user1",
				Password: "Password123!",
			},
			mockSetup: func() {
				mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx *sql.Tx) error) error {
						return fn(nil)
					},
				)
				mockUserGetRepo.EXPECT().GetByParam(gomock.Any(), gomock.Any()).Return(nil, nil)
				mockHasher.EXPECT().Hash("Password123!").Return("", ErrUserRegisterInternal)
			},
			expectedError: ErrUserRegisterInternal,
			expectedResp:  nil,
		},
		{
			name: "Save user failure",
			req: &domain.User{
				Login:    "user1",
				Password: "Password123!",
			},
			mockSetup: func() {
				mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx *sql.Tx) error) error {
						return fn(nil)
					},
				)
				mockUserGetRepo.EXPECT().GetByParam(gomock.Any(), gomock.Any()).Return(nil, nil)
				mockHasher.EXPECT().Hash("Password123!").Return("hashedpassword123", nil)
				mockUserSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(ErrUserRegisterInternal)
			},
			expectedError: ErrUserRegisterInternal,
			expectedResp:  nil,
		},
		{
			name: "JWT encoding failure",
			req: &domain.User{
				Login:    "user1",
				Password: "Password123!",
			},
			mockSetup: func() {
				mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx *sql.Tx) error) error {
						return fn(nil)
					},
				)
				mockUserGetRepo.EXPECT().GetByParam(gomock.Any(), gomock.Any()).Return(nil, nil)
				mockHasher.EXPECT().Hash("Password123!").Return("hashedpassword123", nil)
				mockUserSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
				mockJWTEncoder.EXPECT().Encode("user1").Return("", ErrUserRegisterInternal)
			},
			expectedError: ErrUserRegisterInternal,
			expectedResp:  nil,
		},
		{
			name: "User GetByParam failure",
			req: &domain.User{
				Login:    "user1",
				Password: "Password123!",
			},
			mockSetup: func() {
				mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx *sql.Tx) error) error {
						return fn(nil)
					},
				)
				mockUserGetRepo.EXPECT().GetByParam(gomock.Any(), gomock.Any()).Return(nil, ErrUserRegisterInternal)
				mockHasher.EXPECT().Hash(gomock.Any()).Times(0)
				mockUserSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)
				mockJWTEncoder.EXPECT().Encode(gomock.Any()).Times(0)
			},
			expectedError: ErrUserRegisterInternal,
			expectedResp:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			resp, err := svc.Register(context.Background(), tt.req)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}
		})
	}
}
