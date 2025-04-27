package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

type GophermartUserLoginService struct {
	v          StructValidator
	roExists   ExistsRepository
	roPassword FilterRepository
	pc         PasswordComparer
	jwtGen     JWTGenerator
}

func NewGophermartUserLoginService(
	v StructValidator,
	roExists ExistsRepository,
	roPassword FilterRepository,
	pc PasswordComparer,
	jwtGen JWTGenerator,
) *GophermartUserLoginService {
	return &GophermartUserLoginService{
		v:          v,
		roExists:   roExists,
		roPassword: roPassword,
		pc:         pc,
		jwtGen:     jwtGen,
	}
}

func (svc *GophermartUserLoginService) Login(
	ctx context.Context, req *GophermartUserLoginRequest,
) (*GophermartUserLoginResponse, *APIStatus, *APIStatus) {
	if err := svc.v.Struct(req); err != nil {
		return nil, nil, &APIStatus{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid user login data",
		}
	}
	exists, err := svc.roExists.Exists(ctx, map[string]any{"login": req.Login})
	if err != nil {
		logger.Logger.Errorf("Error checking if user exists: %v", err)
		return nil, nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}
	if !exists {
		return nil, nil, &APIStatus{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid login or password",
		}
	}
	userData, err := svc.roPassword.Filter(ctx, map[string]any{"login": req.Login}, []string{"password"})
	if err != nil {
		logger.Logger.Errorf("Error retrieving password hash: %v", err)
		return nil, nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}
	storedPasswordHash := userData["password"].(string)
	if err := svc.pc.Compare(req.Password, storedPasswordHash); err != nil {
		return nil, nil, &APIStatus{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid login or password",
		}
	}
	token := svc.jwtGen.Generate(map[string]any{"login": req.Login})
	if token == nil {
		return nil, nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error generating authentication token",
		}
	}
	return &GophermartUserLoginResponse{
			Token: *token,
		}, &APIStatus{
			StatusCode: http.StatusOK,
			Message:    "User successfully authenticated",
		}, nil
}

type GophermartUserLoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type GophermartUserLoginResponse struct {
	Token string `json:"token"`
}
