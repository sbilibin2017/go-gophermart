package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserLoginFindRepository interface {
	Find(ctx context.Context, filter *repositories.UserFindFilter) (*repositories.UserFindDB, error)
}

type UserLoginValidator interface {
	Struct(v any) error
}

type UserLoginService struct {
	v                      UserLoginValidator
	repo                   UserLoginFindRepository
	jwtSecretKey           []byte
	generateToken          func(jwtSecretKey []byte, login string) (string, error)
	compareAndHashPassword func(hashedPassword []byte, password []byte) error
}

func NewUserLoginService(
	v UserLoginValidator,
	repo UserLoginFindRepository,
	jwtSecretKey []byte,
	generateToken func(jwtSecretKey []byte, login string) (string, error),
	compareAndHashPassword func(hashedPassword []byte, password []byte) error,
) *UserLoginService {
	return &UserLoginService{
		v:                      v,
		repo:                   repo,
		jwtSecretKey:           jwtSecretKey,
		generateToken:          generateToken,
		compareAndHashPassword: compareAndHashPassword,
	}
}

func (svc *UserLoginService) Login(
	ctx context.Context, req *UserLoginRequest,
) (*UserLoginResponse, *types.APIStatus) {
	if err := svc.v.Struct(req); err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusBadRequest,
			Message: "Invalid request format",
		}
	}

	user, err := svc.repo.Find(ctx, &repositories.UserFindFilter{Login: req.Login})
	if err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusUnauthorized,
			Message: "Invalid login or password",
		}
	}

	if err := svc.compareAndHashPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusUnauthorized,
			Message: "Invalid login or password",
		}
	}

	token, err := svc.generateToken(svc.jwtSecretKey, user.Login)
	if err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusInternalServerError,
			Message: "Failed to generate JWT token",
		}
	}

	return &UserLoginResponse{Token: token}, nil
}

type UserLoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}
