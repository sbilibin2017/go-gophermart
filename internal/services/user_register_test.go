package services

import (
	"context"
	"database/sql"
	e "errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/stretchr/testify/assert"
)

type testMocks struct {
	uow  *MockUnitOfWork
	ugr  *MockUserFilterRepo // Corrected to use the right mock interface
	usr  *MockUserSaveRepo
	hash *MockHasher
	jwt  *MockJWTGenerator
}

func setupTest(t *testing.T) (*UserRegisterService, *testMocks, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	m := &testMocks{
		uow:  NewMockUnitOfWork(ctrl),
		ugr:  NewMockUserFilterRepo(ctrl), // Corrected to use the right mock interface
		usr:  NewMockUserSaveRepo(ctrl),
		hash: NewMockHasher(ctrl),
		jwt:  NewMockJWTGenerator(ctrl),
	}

	svc := NewUserRegisterService(m.ugr, m.usr, m.uow, m.hash, m.jwt)
	return svc, m, ctrl
}

func TestUserRegisterService_Register(t *testing.T) {
	tests := []struct {
		name          string
		user          *User
		setupMocks    func(m *testMocks)
		expectedErr   error
		expectedToken string
	}{
		{
			name: "Success",
			user: &User{Login: "john", Password: "123"},
			setupMocks: func(m *testMocks) {
				m.uow.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(*sql.Tx) error) error {
						m.ugr.EXPECT().Filter(ctx, &repositories.UserFilter{Login: "john"}).Return(nil, nil)
						m.hash.EXPECT().Hash("123").Return("hashed123")
						m.usr.EXPECT().Save(ctx, &repositories.UserSave{Login: "john", Password: "hashed123"}).Return(nil)
						m.jwt.EXPECT().Generate("john").Return("token123")
						return fn(nil)
					},
				)
			},
			expectedErr:   nil,
			expectedToken: "token123",
		},
		{
			name: "UserAlreadyExists",
			user: &User{Login: "john", Password: "123"},
			setupMocks: func(m *testMocks) {
				m.uow.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(*sql.Tx) error) error {
						m.ugr.EXPECT().Filter(ctx, &repositories.UserFilter{Login: "john"}).Return(&repositories.UserFiltered{}, nil)
						return fn(nil)
					},
				)
			},
			expectedErr:   ErrUserAlreadyExists,
			expectedToken: "",
		},
		{
			name: "GetUserError",
			user: &User{Login: "john", Password: "123"},
			setupMocks: func(m *testMocks) {
				m.uow.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(*sql.Tx) error) error {
						m.ugr.EXPECT().Filter(ctx, &repositories.UserFilter{Login: "john"}).Return(nil, e.New("db error"))
						return fn(nil)
					},
				)
			},
			expectedErr:   ErrUserRegisterInternal,
			expectedToken: "",
		},
		{
			name: "SaveError",
			user: &User{Login: "john", Password: "123"},
			setupMocks: func(m *testMocks) {
				m.uow.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(*sql.Tx) error) error {
						m.ugr.EXPECT().Filter(ctx, &repositories.UserFilter{Login: "john"}).Return(nil, nil)
						m.hash.EXPECT().Hash("123").Return("hashed123")
						m.usr.EXPECT().Save(ctx, &repositories.UserSave{Login: "john", Password: "hashed123"}).Return(e.New("insert error"))
						return fn(nil)
					},
				)
			},
			expectedErr:   ErrUserRegisterInternal,
			expectedToken: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, m, ctrl := setupTest(t)
			defer ctrl.Finish()

			tt.setupMocks(m)

			token, err := svc.Register(context.Background(), tt.user)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedToken, token)
		})
	}
}
