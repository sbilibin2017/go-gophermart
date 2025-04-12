package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sbilibin2017/go-gophermart/internal/repositories"
)

type UserFilterRepo interface {
	Filter(ctx context.Context, filter *repositories.UserFilter) (*repositories.UserFiltered, error)
}

type UserSaveRepo interface {
	Save(ctx context.Context, user *repositories.UserSave) error
}

type Hasher interface {
	Hash(password string) string
}

type JWTGenerator interface {
	Generate(login string) string
}

type UnitOfWork interface {
	Do(ctx context.Context, operation func(tx *sql.Tx) error) error
}

type UserRegisterService struct {
	ugr UserFilterRepo
	usr UserSaveRepo
	h   Hasher
	g   JWTGenerator
	uow UnitOfWork
}

func NewUserRegisterService(
	ugr UserFilterRepo,
	usr UserSaveRepo,
	h Hasher,
	g JWTGenerator,
	uow UnitOfWork,
) *UserRegisterService {
	return &UserRegisterService{
		ugr: ugr,
		usr: usr,
		h:   h,
		g:   g,
		uow: uow,
	}
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (svc *UserRegisterService) Register(ctx context.Context, u *User) (string, error) {
	var token string
	err := svc.uow.Do(ctx, func(tx *sql.Tx) error {
		user, err := svc.ugr.Filter(ctx, &repositories.UserFilter{Login: u.Login})
		if err != nil {
			return err
		}
		if user != nil {
			return ErrUserAlreadyExists
		}

		u.Password = svc.h.Hash(u.Password)

		err = svc.usr.Save(ctx, &repositories.UserSave{
			Login:    u.Login,
			Password: u.Password,
		})
		if err != nil {
			return err
		}

		token = svc.g.Generate(u.Login)
		return nil
	})

	if err != nil {
		return "", err
	}

	return token, nil
}

var (
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrUserRegisterInternal = errors.New("internal error")
)
