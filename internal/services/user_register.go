package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/jwt"
	"github.com/sbilibin2017/go-gophermart/internal/password"
)

type UserGetRepo interface {
	GetByParam(ctx context.Context, p map[string]any) (map[string]any, error)
}

type UserSaveRepo interface {
	Save(ctx context.Context, u map[string]any) error
}

type UnitOfWork interface {
	Do(ctx context.Context, operation func(tx *sql.Tx) error) error
}

type UserRegisterService struct {
	config *configs.GophermartConfig
	ugr    UserGetRepo
	usr    UserSaveRepo
	uow    UnitOfWork
}

func NewUserRegisterService(
	config *configs.GophermartConfig,
	ugr UserGetRepo,
	usr UserSaveRepo,
	uow UnitOfWork,
) *UserRegisterService {
	return &UserRegisterService{
		config: config,
		ugr:    ugr,
		usr:    usr,
		uow:    uow,
	}
}

var (
	ErrUserRegisterInternal = errors.New("internal error")
	ErrUserAlreadyExists    = errors.New("user already exists")
)

func (svc *UserRegisterService) Register(
	ctx context.Context, u *domain.User,
) (string, error) {
	var token string
	err := svc.uow.Do(ctx, func(tx *sql.Tx) error {
		data := map[string]any{
			"login": u.Login,
		}
		user, err := svc.ugr.GetByParam(ctx, data)
		if err != nil {
			return ErrUserRegisterInternal
		}
		if user != nil {
			return ErrUserAlreadyExists
		}
		hashedPassword, err := password.Hash(u.Password)
		if err != nil {
			return ErrUserRegisterInternal
		}
		data["password"] = hashedPassword
		err = svc.usr.Save(ctx, data)
		if err != nil {
			return ErrUserRegisterInternal
		}
		token, err = jwt.Generate(u.Login, svc.config.JWTSecretKey, svc.config.JWTExp)
		if err != nil {
			return ErrUserRegisterInternal
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return token, nil
}
