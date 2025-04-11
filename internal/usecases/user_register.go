package usecases

import (
	"context"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
)

type LoginValidator interface {
	Validate(login string) error
}

type PasswordValidator interface {
	Validate(password string) error
}

type UserRegisterService interface {
	Register(ctx context.Context, u *domain.User) (string, error)
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
}

func (uc *UserRegisterUsecase) Execute(
	ctx context.Context, req *UserRegisterRequest,
) (*UserRegisterResponse, error) {
	err := uc.lv.Validate(req.Login)
	if err != nil {
		return nil, err
	}
	err = uc.pv.Validate(req.Password)
	if err != nil {
		return nil, err
	}
	token, err := uc.svc.Register(ctx, &domain.User{
		Login: req.Login, Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &UserRegisterResponse{AccessToken: token}, nil
}
