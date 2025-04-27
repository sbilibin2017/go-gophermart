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

type GophermartUserRegisterExistsRepository interface {
	Exists(ctx context.Context, login string) (bool, error)
}

type GophermartUserRegisterSaveRepository interface {
	Save(ctx context.Context, login, password string) error
}

type GophermartUserRegisterService struct {
	v         GophermartUserRegisterValidator
	ro        GophermartUserRegisterExistsRepository
	rs        GophermartUserRegisterSaveRepository
	jwtSecret string
	jwtExp    time.Duration
}

type GophermartUserRegisterValidator interface {
	Struct(v any) error
}

func NewGophermartUserRegisterService(
	v GophermartUserRegisterValidator,
	ro GophermartUserRegisterExistsRepository,
	rs GophermartUserRegisterSaveRepository,
	jwtSecret string,
	jwtExp time.Duration,
) *GophermartUserRegisterService {
	return &GophermartUserRegisterService{
		v:         v,
		ro:        ro,
		rs:        rs,
		jwtSecret: jwtSecret,
		jwtExp:    jwtExp,
	}
}

func (svc *GophermartUserRegisterService) Register(
	ctx context.Context, req *types.GophermartUserRegisterRequest,
) (*types.GophermartUserRegisterResponse, *types.APIStatus, *types.APIStatus) {
	if err := svc.v.Struct(req); err != nil {
		return nil, nil, &types.APIStatus{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid user registration data",
		}
	}
	exists, err := svc.ro.Exists(ctx, req.Login)
	if err != nil {
		logger.Logger.Errorf("Error checking if user exists: %v", err)
		return nil, nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}
	if exists {
		return nil, nil, &types.APIStatus{
			StatusCode: http.StatusConflict,
			Message:    "Login already taken",
		}
	}
	passwordHash, err := password.Hash(req.Password)
	if err != nil {
		return nil, nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error hashing password",
		}
	}
	if err := svc.rs.Save(ctx, req.Login, passwordHash); err != nil {
		logger.Logger.Errorf("Error saving user: %v", err)
		return nil, nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}
	token, err := jwt.GenerateToken(req.Login, svc.jwtSecret, 24*time.Hour)
	if err != nil {
		return nil, nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error generating authentication token",
		}
	}

	return &types.GophermartUserRegisterResponse{
			Token: token,
		}, &types.APIStatus{
			StatusCode: http.StatusOK,
			Message:    "User successfully registered and authenticated",
		}, nil
}
