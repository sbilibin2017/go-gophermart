package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
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
	Hash(password string) (string, error)
}

type JWTEncoder interface {
	Encode(login string) (string, error)
}

type UserRegisterService struct {
	ugr UserGetRepo
	usr UserSaveRepo
	uow UnitOfWork
	h   Hasher
	e   JWTEncoder
}

func NewUserRegisterService(
	ugr UserGetRepo,
	usr UserSaveRepo,
	uow UnitOfWork,
	h Hasher,
	e JWTEncoder,
) *UserRegisterService {
	return &UserRegisterService{
		ugr: ugr,
		usr: usr,
		uow: uow,
		h:   h,
		e:   e,
	}
}

var (
	ErrUserRegisterInternal = errors.New("internal error")
	ErrUserAlreadyExists    = errors.New("user already exists")
)

func (svc *UserRegisterService) Register(
	ctx context.Context, u *domain.User,
) (*domain.UserToken, error) {
	var userToken domain.UserToken
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
		hashedPassword, err := svc.h.Hash(u.Password)
		if err != nil {
			return ErrUserRegisterInternal
		}
		data["password"] = hashedPassword
		err = svc.usr.Save(ctx, data)
		if err != nil {
			return ErrUserRegisterInternal
		}
		token, err := svc.e.Encode(u.Login)
		if err != nil {
			return ErrUserRegisterInternal
		}
		userToken.Access = token
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &userToken, nil
}
