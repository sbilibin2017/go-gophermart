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

func TestUserRegisterService_Register_Success(t *testing.T) {
	svc, m, ctrl := setupTest(t)
	defer ctrl.Finish()

	u := &User{Login: "john", Password: "123"}

	m.uow.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(*sql.Tx) error) error {
			m.ugr.EXPECT().Filter(ctx, &repositories.UserFilter{Login: "john"}).Return(nil, nil) // Corrected method call
			m.hash.EXPECT().Hash("123").Return("hashed123")
			m.usr.EXPECT().Save(ctx, &repositories.UserSave{Login: "john", Password: "hashed123"}).Return(nil) // Corrected method call
			m.jwt.EXPECT().Generate("john").Return("token123")
			return fn(nil)
		},
	)

	token, err := svc.Register(context.Background(), u)
	assert.NoError(t, err)
	assert.Equal(t, "token123", token)
}

func TestUserRegisterService_Register_UserAlreadyExists(t *testing.T) {
	svc, m, ctrl := setupTest(t)
	defer ctrl.Finish()

	u := &User{Login: "john", Password: "123"}

	m.uow.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(*sql.Tx) error) error {
			m.ugr.EXPECT().Filter(ctx, &repositories.UserFilter{Login: "john"}).Return(&repositories.UserFiltered{}, nil) // Corrected method call
			return fn(nil)
		},
	)

	token, err := svc.Register(context.Background(), u)
	assert.ErrorIs(t, err, ErrUserAlreadyExists)
	assert.Empty(t, token)
}

func TestUserRegisterService_Register_GetUserError(t *testing.T) {
	svc, m, ctrl := setupTest(t)
	defer ctrl.Finish()

	u := &User{Login: "john", Password: "123"}

	m.uow.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(*sql.Tx) error) error {
			m.ugr.EXPECT().Filter(ctx, &repositories.UserFilter{Login: "john"}).Return(nil, e.New("db error")) // Corrected method call
			return fn(nil)
		},
	)

	token, err := svc.Register(context.Background(), u)
	assert.ErrorIs(t, err, ErrUserRegisterInternal)
	assert.Empty(t, token)
}

func TestUserRegisterService_Register_SaveError(t *testing.T) {
	svc, m, ctrl := setupTest(t)
	defer ctrl.Finish()

	u := &User{Login: "john", Password: "123"}

	m.uow.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(*sql.Tx) error) error {
			m.ugr.EXPECT().Filter(ctx, &repositories.UserFilter{Login: "john"}).Return(nil, nil) // Corrected method call
			m.hash.EXPECT().Hash("123").Return("hashed123")
			m.usr.EXPECT().Save(ctx, &repositories.UserSave{Login: "john", Password: "hashed123"}).Return(e.New("insert error"))
			return fn(nil)
		},
	)

	token, err := svc.Register(context.Background(), u)
	assert.ErrorIs(t, err, ErrUserRegisterInternal)
	assert.Empty(t, token)
}
