package services

import (
	"context"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/errors"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
)

type UserGetRepo interface {
	GetByParam(ctx context.Context, p *repositories.UserGetParam) (*repositories.UserGet, error)
}

type UserSaveRepo interface {
	Save(ctx context.Context, u *repositories.UserSave) error
}

type UserService struct {
	ugr UserGetRepo
	usr UserSaveRepo
}

func (svc *UserService) Register(ctx context.Context, u *domain.User) error {
	user, err := svc.ugr.GetByParam(ctx, &repositories.UserGetParam{Login: u.Login})
	if err != nil {
		return errors.ErrInternal
	}
	if user != nil {
		return errors.ErrUserAlreadyExists
	}
	err = svc.usr.Save(
		ctx, &repositories.UserSave{Login: u.Login, Password: u.Password},
	)
	if err != nil {
		return errors.ErrInternal
	}
	return nil
}
