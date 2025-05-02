package services

import (
	"context"
	"net/http"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserRegisterUserFilterOneRepository interface {
	FilterOne(ctx context.Context, login string) (*types.UserDB, error)
}

type UserRegisterUserSaveRepository interface {
	Save(ctx context.Context, login string, password string) error
}

type UserRegisterValidator interface {
	Struct(v any) error
}

type UserRegisterService struct {
	v            UserRegisterValidator
	uf           UserRegisterUserFilterOneRepository
	us           UserRegisterUserSaveRepository
	jwtSecretKey string
	jwtExp       time.Duration
	issuer       string
}

func NewUserRegisterService(
	v UserRegisterValidator,
	uf UserRegisterUserFilterOneRepository,
	us UserRegisterUserSaveRepository,
	jwtSecretKey string,
	jwtExp time.Duration,
	issuer string,
) *UserRegisterService {
	return &UserRegisterService{
		v:            v,
		uf:           uf,
		us:           us,
		jwtSecretKey: jwtSecretKey,
		jwtExp:       jwtExp,
		issuer:       issuer,
	}
}

func (svc *UserRegisterService) Register(
	ctx context.Context, req *types.UserRegisterRequest,
) (string, *types.APIStatus, *types.APIStatus) {
	if err := svc.v.Struct(req); err != nil {
		valErr := formatValidationError(err)
		if valErr != nil {
			return "", nil, &types.APIStatus{
				StatusCode: http.StatusBadRequest,
				Message:    valErr.Message,
			}
		}
	}

	existingUser, err := svc.uf.FilterOne(ctx, req.Login)
	if err != nil {
		return "", nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error while registering user",
		}
	}
	if existingUser != nil {
		return "", nil, &types.APIStatus{
			StatusCode: http.StatusConflict,
			Message:    "User login already exists",
		}
	}

	passwordHashed, err := hashPassword(req.Password)
	if err != nil {
		return "", nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error while registering user",
		}
	}

	if err := svc.us.Save(ctx, req.Login, passwordHashed); err != nil {
		return "", nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error while registering user",
		}
	}

	token, err := generateTokenString(
		svc.jwtSecretKey,
		svc.jwtExp,
		svc.issuer,
		req.Login,
	)
	if err != nil {
		return "", nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error while registering user",
		}
	}

	return token, &types.APIStatus{
		StatusCode: http.StatusCreated,
		Message:    "User successfully registered",
	}, nil
}
