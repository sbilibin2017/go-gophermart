package services

import (
	"context"
	"database/sql"
	e "errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type testMocks struct {
	uow  *MockUnitOfWork
	ugr  *MockUserGetRepo
	usr  *MockUserSaveRepo
	hash *MockHasher
	jwt  *MockJWTGenerator
}

func setupTest(t *testing.T) (*UserRegisterService, *testMocks, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	m := &testMocks{
		uow:  NewMockUnitOfWork(ctrl),
		ugr:  NewMockUserGetRepo(ctrl),
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
			m.ugr.EXPECT().GetByParam(ctx, map[string]any{"login": "john"}).Return(nil, nil)
			m.hash.EXPECT().Hash("123").Return("hashed123")
			m.usr.EXPECT().Save(ctx, map[string]any{"login": "john", "password": "hashed123"}).Return(nil)
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
			m.ugr.EXPECT().GetByParam(ctx, map[string]any{"login": "john"}).Return(map[string]any{"id": 1}, nil)
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
			m.ugr.EXPECT().GetByParam(ctx, map[string]any{"login": "john"}).Return(nil, e.New("db error"))
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
			m.ugr.EXPECT().GetByParam(ctx, map[string]any{"login": "john"}).Return(nil, nil)
			m.hash.EXPECT().Hash("123").Return("hashed123")
			m.usr.EXPECT().Save(ctx, map[string]any{"login": "john", "password": "hashed123"}).Return(e.New("insert error"))
			return fn(nil)
		},
	)

	token, err := svc.Register(context.Background(), u)
	assert.ErrorIs(t, err, ErrUserRegisterInternal)
	assert.Empty(t, token)
}
