package usecases

import (
	"context"
	"database/sql"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/requests"
	"github.com/sbilibin2017/go-gophermart/internal/responses"
)

type UserService interface {
	Register(ctx context.Context, u *domain.User) error
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

type UserRegisterUsecase struct {
	svc UserService
	uow UnitOfWork
	h   Hasher
	e   JWTEncoder
}

func (uc *UserRegisterUsecase) Execute(
	ctx context.Context, req *requests.UserRegisterRequest,
) (*responses.UserRegisterResponse, error) {
	var token string
	err := uc.uow.Do(ctx, func(tx *sql.Tx) error {
		if err := req.Validate(); err != nil {
			return err
		}
		hashedPassword, err := uc.h.Hash(req.Password)
		if err != nil {
			return err
		}
		user := &domain.User{
			Login:    req.Login,
			Password: hashedPassword,
		}
		if err := uc.svc.Register(ctx, user); err != nil {
			return err
		}
		tok, err := uc.e.Encode(user.Login)
		if err != nil {
			return err
		}
		token = tok

		return nil
	})
	if err != nil {
		return nil, err
	}
	return &responses.UserRegisterResponse{AccessToken: token}, nil
}
