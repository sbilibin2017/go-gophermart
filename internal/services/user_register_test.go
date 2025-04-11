package services

import (
	"context"
	"database/sql"
	e "errors"

	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userGetRepo := NewMockUserGetRepo(ctrl)
	userSaveRepo := NewMockUserSaveRepo(ctrl)
	unitOfWork := NewMockUnitOfWork(ctrl)
	hasher := NewMockHasher(ctrl)
	jwtGen := NewMockJWTGenerator(ctrl)

	ctx := context.Background()
	user := &domain.User{
		Login:    "testuser",
		Password: "1234",
	}

	hashedPassword := "hashed1234"

	// Настроим mocks
	userGetRepo.EXPECT().
		GetByParam(ctx, map[string]any{"login": "testuser"}).
		Return(nil, nil)

	hasher.EXPECT().
		Hash("1234").
		Return(hashedPassword)

	userSaveRepo.EXPECT().
		Save(ctx, map[string]any{
			"login":    "testuser",
			"password": hashedPassword,
		}).
		Return(nil)

	unitOfWork.EXPECT().
		Do(ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, op func(tx *sql.Tx) error) error {
			return op(nil)
		})

	jwtGen.EXPECT().
		Generate("testuser").
		Return("jwt_token")

	// Подготовка конфига
	cfg := &configs.GophermartConfig{}

	service := NewUserRegisterService(cfg, userGetRepo, userSaveRepo, unitOfWork, hasher, jwtGen)

	// Выполнение
	token, err := service.Register(ctx, user)

	// Проверки
	require.NoError(t, err)
	assert.Equal(t, "jwt_token", token)
}

func TestRegisterUser_GetByParamError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userGetRepo := NewMockUserGetRepo(ctrl)
	userSaveRepo := NewMockUserSaveRepo(ctrl)
	unitOfWork := NewMockUnitOfWork(ctrl)
	hasher := NewMockHasher(ctrl)
	jwtGen := NewMockJWTGenerator(ctrl)
	cfg := &configs.GophermartConfig{}

	service := NewUserRegisterService(cfg, userGetRepo, userSaveRepo, unitOfWork, hasher, jwtGen)

	ctx := context.Background()
	user := &domain.User{Login: "testuser", Password: "1234"}

	userGetRepo.EXPECT().
		GetByParam(ctx, map[string]any{"login": "testuser"}).
		Return(nil, e.New("db error"))

	unitOfWork.EXPECT().
		Do(ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, op func(tx *sql.Tx) error) error {
			return op(nil)
		})

	token, err := service.Register(ctx, user)

	require.Error(t, err)
	assert.Equal(t, "", token)
	assert.ErrorIs(t, err, errors.ErrInternal)
}

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userGetRepo := NewMockUserGetRepo(ctrl)
	userSaveRepo := NewMockUserSaveRepo(ctrl)
	unitOfWork := NewMockUnitOfWork(ctrl)
	hasher := NewMockHasher(ctrl)
	jwtGen := NewMockJWTGenerator(ctrl)
	cfg := &configs.GophermartConfig{}

	service := NewUserRegisterService(cfg, userGetRepo, userSaveRepo, unitOfWork, hasher, jwtGen)

	ctx := context.Background()
	user := &domain.User{Login: "testuser", Password: "1234"}

	userGetRepo.EXPECT().
		GetByParam(ctx, map[string]any{"login": "testuser"}).
		Return(map[string]any{"login": "testuser"}, nil)

	unitOfWork.EXPECT().
		Do(ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, op func(tx *sql.Tx) error) error {
			return op(nil)
		})

	token, err := service.Register(ctx, user)

	require.Error(t, err)
	assert.Equal(t, "", token)
	assert.ErrorIs(t, err, errors.ErrUserAlreadyExists)
}

func TestRegisterUser_SaveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userGetRepo := NewMockUserGetRepo(ctrl)
	userSaveRepo := NewMockUserSaveRepo(ctrl)
	unitOfWork := NewMockUnitOfWork(ctrl)
	hasher := NewMockHasher(ctrl)
	jwtGen := NewMockJWTGenerator(ctrl)
	cfg := &configs.GophermartConfig{}

	service := NewUserRegisterService(cfg, userGetRepo, userSaveRepo, unitOfWork, hasher, jwtGen)

	ctx := context.Background()
	user := &domain.User{Login: "testuser", Password: "1234"}
	hashedPassword := "hashed1234"

	userGetRepo.EXPECT().
		GetByParam(ctx, map[string]any{"login": "testuser"}).
		Return(nil, nil)

	hasher.EXPECT().
		Hash("1234").
		Return(hashedPassword)

	userSaveRepo.EXPECT().
		Save(ctx, map[string]any{
			"login":    "testuser",
			"password": hashedPassword,
		}).
		Return(e.New("insert error"))

	unitOfWork.EXPECT().
		Do(ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, op func(tx *sql.Tx) error) error {
			return op(nil)
		})

	token, err := service.Register(ctx, user)

	require.Error(t, err)
	assert.Equal(t, "", token)
	assert.ErrorIs(t, err, errors.ErrInternal)
}
