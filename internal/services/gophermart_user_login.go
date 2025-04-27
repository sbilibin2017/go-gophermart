package services

import (
	"context"
	"net/http"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/jwt"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/password"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type GophermartUserLoginExistsRepository interface {
	ExistsByLogin(ctx context.Context, login string) (bool, error)
}

type GophermartUserLoginPasswordRepository interface {
	GetPasswordByLogin(ctx context.Context, login string) (string, error)
}

type GophermartUserLoginService interface {
	Login(ctx context.Context, req *types.GophermartUserLoginRequest) (*types.GophermartUserLoginResponse, *types.APIStatus, *types.APIStatus)
}

type GophermartUserLoginServiceImpl struct {
	roExists   GophermartUserLoginExistsRepository
	roPassword GophermartUserLoginPasswordRepository
	jwtSecret  string
	jwtExp     time.Duration
}

func NewGophermartUserLoginService(
	roExists GophermartUserLoginExistsRepository,
	roPassword GophermartUserLoginPasswordRepository,
	jwtSecret string,
	jwtExp time.Duration,
) *GophermartUserLoginServiceImpl {
	return &GophermartUserLoginServiceImpl{
		roExists:   roExists,
		roPassword: roPassword,
		jwtSecret:  jwtSecret,
		jwtExp:     jwtExp,
	}
}

func (svc *GophermartUserLoginServiceImpl) Login(
	ctx context.Context, req *types.GophermartUserLoginRequest,
) (*types.GophermartUserLoginResponse, *types.APIStatus, *types.APIStatus) {
	exists, err := svc.roExists.ExistsByLogin(ctx, req.Login)
	if err != nil {
		logger.Logger.Errorf("Error checking if user exists: %v", err)
		return nil, nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}
	if !exists {
		return nil, nil, &types.APIStatus{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid login or password",
		}
	}
	storedPasswordHash, err := svc.roPassword.GetPasswordByLogin(ctx, req.Login)
	if err != nil {
		logger.Logger.Errorf("Error retrieving password hash: %v", err)
		return nil, nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}
	if err := password.Compare(req.Password, storedPasswordHash); err != nil {
		return nil, nil, &types.APIStatus{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid login or password",
		}
	}
	token, err := jwt.GenerateToken(req.Login, svc.jwtSecret, svc.jwtExp)
	if err != nil {
		return nil, nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error generating authentication token",
		}
	}
	return &types.GophermartUserLoginResponse{
			Token: token,
		}, &types.APIStatus{
			StatusCode: http.StatusOK,
			Message:    "User successfully authenticated",
		}, nil
}
