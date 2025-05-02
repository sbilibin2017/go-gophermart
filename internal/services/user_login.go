package services

import (
	"context"
	"net/http"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserLoginUserFilterOneRepository interface {
	FilterOne(ctx context.Context, login string) (*types.UserDB, error)
}

type UserLoginValidator interface {
	Struct(v any) error
}

type UserLoginService struct {
	v            UserLoginValidator
	uf           UserLoginUserFilterOneRepository
	jwtSecretKey string
	jwtExp       time.Duration
	issuer       string
}

func NewUserLoginService(
	v UserLoginValidator,
	uf UserLoginUserFilterOneRepository,
	jwtSecretKey string,
	jwtExp time.Duration,
	issuer string,
) *UserLoginService {
	return &UserLoginService{
		v:            v,
		uf:           uf,
		jwtSecretKey: jwtSecretKey,
		jwtExp:       jwtExp,
		issuer:       issuer,
	}
}

func (svc *UserLoginService) Login(
	ctx context.Context, req *types.UserLoginRequest,
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

	err = comparePassword(existingUser.Password, req.Password)
	if err != nil {
		return "", nil, &types.APIStatus{
			StatusCode: http.StatusConflict,
			Message:    "Incorrect user password",
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
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid request",
	}, nil
}
