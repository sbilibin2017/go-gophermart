package services

import (
	"context"
	"database/sql"

	"time"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/errors"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserGetRepo interface {
	GetByParam(ctx context.Context, p *types.UserGetParam) (*types.User, error)
}

type UserSaveRepo interface {
	Save(ctx context.Context, u *types.User) error
}

type UnitOfWork interface {
	Do(ctx context.Context, operation func(tx *sql.Tx) error) error
}

type UserService struct {
	config         *configs.GophermartConfig
	ugr            UserGetRepo
	usr            UserSaveRepo
	uow            UnitOfWork
	tokenEncoder   func(secretKey string, exp time.Duration, login string) (string, error)
	loginValidator func(login string) error
	hasher         func(password string) (string, error)
}

func (svc *UserService) Register(
	ctx context.Context,
	u *types.User,
) (*types.Token, error) {
	var token string
	err := svc.uow.Do(ctx, func(tx *sql.Tx) error {
		err := svc.loginValidator(u.Login)
		if err != nil {
			return err
		}
		user, err := svc.ugr.GetByParam(ctx, &types.UserGetParam{Login: u.Login})
		if err != nil {
			return errors.ErrInternal
		}
		if user != nil {
			return errors.ErrUserAlreadyExists
		}
		passwordHash, err := svc.hasher(u.Password)
		if err != nil {
			return errors.ErrInternal
		}
		u.Password = passwordHash
		err = svc.usr.Save(ctx, u)
		if err != nil {
			return errors.ErrInternal
		}
		token, err = svc.tokenEncoder(svc.config.JWTSecretKey, svc.config.JWTExp, u.Login)
		if err != nil {
			return errors.ErrInternal
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &types.Token{Access: token}, nil
}
