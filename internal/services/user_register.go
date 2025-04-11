package services

import (
	"context"
	"database/sql"
	"errors"
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
	ugr UserGetRepo
	usr UserSaveRepo
	uow UnitOfWork
	h   Hasher
	g   JWTGenerator
}

func NewUserRegisterService(
	ugr UserGetRepo,
	usr UserSaveRepo,
	uow UnitOfWork,
	h Hasher,
	g JWTGenerator,
) *UserRegisterService {
	return &UserRegisterService{
		ugr: ugr,
		usr: usr,
		uow: uow,
		h:   h,
		g:   g,
	}
}

type User struct {
	Login    string
	Password string
}

func (svc *UserRegisterService) Register(
	ctx context.Context, u *User,
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

		data["password"] = svc.h.Hash(u.Password)
		err = svc.usr.Save(ctx, data)
		if err != nil {
			return ErrUserRegisterInternal
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
