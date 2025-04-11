package services

import (
	"context"
	"database/sql"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/errors"
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

type Hasher interface {
	Hash(password string) string
}

type JWTGenerator interface {
	Generate(login string) string
}

type UserRegisterService struct {
	config *configs.GophermartConfig
	ugr    UserGetRepo
	usr    UserSaveRepo
	uow    UnitOfWork
	h      Hasher
	g      JWTGenerator
}

func NewUserRegisterService(
	config *configs.GophermartConfig,
	ugr UserGetRepo,
	usr UserSaveRepo,
	uow UnitOfWork,
	h Hasher,
	g JWTGenerator,
) *UserRegisterService {
	return &UserRegisterService{
		config: config,
		ugr:    ugr,
		usr:    usr,
		uow:    uow,
		h:      h,
		g:      g,
	}
}

var ()

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
			return errors.ErrInternal
		}
		if user != nil {
			return errors.ErrUserAlreadyExists
		}
		data["password"] = svc.h.Hash(u.Password)
		err = svc.usr.Save(ctx, data)
		if err != nil {
			return errors.ErrInternal
		}
		token = svc.g.Generate(u.Login)
		return nil
	})
	if err != nil {
		return "", err
	}
	return token, nil
}
