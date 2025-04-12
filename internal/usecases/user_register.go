package usecases

import (
	"context"
	"database/sql"

	"github.com/sbilibin2017/go-gophermart/internal/services"
)

type LoginValidator interface {
	Validate(login string) error
}

type PasswordValidator interface {
	Validate(password string) error
}

type UserRegisterService interface {
	Register(ctx context.Context, u *services.User) (string, error)
}

type UnitOfWork interface {
	Do(ctx context.Context, operation func(tx *sql.Tx) error) error
}

type UserRegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserRegisterResponse struct {
	AccessToken string `json:"access_token"`
}

type UserRegisterUsecase struct {
	lv  LoginValidator
	pv  PasswordValidator
	svc UserRegisterService
	uow UnitOfWork
}

func NewUserRegisterUsecase(
	uow UnitOfWork,
	lv LoginValidator,
	pv PasswordValidator,
	svc UserRegisterService,

) *UserRegisterUsecase {
	return &UserRegisterUsecase{
		uow: uow,
		lv:  lv,
		pv:  pv,
		svc: svc,
	}
}

func (uc *UserRegisterUsecase) Execute(
	ctx context.Context, req *UserRegisterRequest,
) (*UserRegisterResponse, error) {
	var token string

	err := uc.uow.Do(ctx, func(tx *sql.Tx) error {
		err := uc.lv.Validate(req.Login)
		if err != nil {
			return err
		}
		err = uc.pv.Validate(req.Password)
		if err != nil {
			return err
		}
		token, err = uc.svc.Register(ctx, &services.User{
			Login: req.Login, Password: req.Password,
		})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &UserRegisterResponse{AccessToken: token}, nil
}
