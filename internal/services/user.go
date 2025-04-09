package services

import (
	"context"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/errors"
)

type UserGetRepo interface {
	GetByParam(ctx context.Context, p *domain.UserGetParam) (*domain.User, error)
}

type UserSaveRepo interface {
	Save(ctx context.Context, u *domain.User) error
}

type UserService struct {
	ugr UserGetRepo
	usr UserSaveRepo
}

func (svc *UserService) Register(ctx context.Context, u *domain.User) error {
	user, err := svc.ugr.GetByParam(ctx, &domain.UserGetParam{Login: u.Login})
	if err != nil {
		return errors.ErrInternal
	}
	if user != nil {
		return errors.ErrUserAlreadyExists
	}
	err = svc.usr.Save(ctx, u)
	if err != nil {
		return errors.ErrInternal
	}
	return nil
}
