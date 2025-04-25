package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/types"
	"golang.org/x/crypto/bcrypt"
)

type UserRegisterExistsRepository interface {
	Exists(ctx context.Context, login *repositories.UserFindFilter) (bool, error)
}

type UserRegisterSaveRepository interface {
	Save(ctx context.Context, user *repositories.UserSave) error
}

type UserRegisterValidator interface {
	Struct(v any) error
}

type UserRegisterService struct {
	v                        UserRegisterValidator
	ure                      UserRegisterExistsRepository
	urs                      UserRegisterSaveRepository
	jwtSecretKey             []byte
	generateToken            func(jwtSecretKey []byte, login string) (string, error)
	generateHashFromPassword func(password []byte, cost int) ([]byte, error)
}

func NewUserRegisterService(
	v OrderAcceptValidator,
	ure UserRegisterExistsRepository,
	urs UserRegisterSaveRepository,
	jwtSecretKey []byte,
	generateToken func(jwtSecretKey []byte, login string) (string, error),
	generateHashFromPassword func(password []byte, cost int) ([]byte, error),
) *UserRegisterService {
	return &UserRegisterService{
		v:                        v,
		ure:                      ure,
		urs:                      urs,
		jwtSecretKey:             jwtSecretKey,
		generateToken:            generateToken,
		generateHashFromPassword: generateHashFromPassword,
	}
}

func (svc *UserRegisterService) Register(
	ctx context.Context, req *UserRegisterRequest,
) (*UserRegisterResponse, *types.APIStatus) {
	if err := svc.v.Struct(req); err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusBadRequest,
			Message: "Invalid request format",
		}
	}

	exists, err := svc.ure.Exists(ctx, &repositories.UserFindFilter{Login: req.Login})
	if err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}
	if exists {
		return nil, &types.APIStatus{
			Status:  http.StatusConflict,
			Message: "Login already taken",
		}
	}

	hashedPassword, err := svc.generateHashFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusInternalServerError,
			Message: "Error while hashing password",
		}
	}

	user := &repositories.UserSave{
		Login:    req.Login,
		Password: string(hashedPassword),
	}

	if err := svc.urs.Save(ctx, user); err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}

	token, err := svc.generateToken(svc.jwtSecretKey, user.Login)
	if err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusInternalServerError,
			Message: "Failed to generate JWT token",
		}
	}

	return &UserRegisterResponse{Token: token}, nil
}

type UserRegisterRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserRegisterResponse struct {
	Token string `json:"token"`
}
