package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

type GophermartUserRegisterService struct {
	v      StructValidator
	ro     ExistsRepository
	rs     SaveRepository
	ph     PasswordHasher
	jwtGen JWTGenerator
}

func NewGophermartUserRegisterService(
	v StructValidator,
	ro ExistsRepository,
	rs SaveRepository,
	ph PasswordHasher,
	jwtGen JWTGenerator,
) *GophermartUserRegisterService {
	return &GophermartUserRegisterService{
		v:      v,
		ro:     ro,
		rs:     rs,
		ph:     ph,
		jwtGen: jwtGen,
	}
}

func (svc *GophermartUserRegisterService) Register(
	ctx context.Context, req *GophermartUserRegisterRequest,
) (*GophermartUserRegisterResponse, *APIStatus, *APIStatus) {
	if err := svc.v.Struct(req); err != nil {
		return nil, nil, &APIStatus{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid user registration data",
		}
	}
	// Используем интерфейс ExistsRepository для проверки существования пользователя
	exists, err := svc.ro.Exists(ctx, map[string]any{"login": req.Login})
	if err != nil {
		logger.Logger.Errorf("Error checking if user exists: %v", err)
		return nil, nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}
	if exists {
		return nil, nil, &APIStatus{
			StatusCode: http.StatusConflict,
			Message:    "Login already taken",
		}
	}

	// Используем интерфейс PasswordHasher для хеширования пароля
	passwordHash := svc.ph.Hash(req.Password)
	if passwordHash == nil {
		return nil, nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error hashing password",
		}
	}

	// Используем интерфейс SaveRepository для сохранения данных пользователя
	if err := svc.rs.Save(ctx, map[string]any{"login": req.Login, "password_hash": *passwordHash}); err != nil {
		logger.Logger.Errorf("Error saving user: %v", err)
		return nil, nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}

	// Генерация JWT токена с использованием интерфейса JWTGenerator
	token := svc.jwtGen.Generate(map[string]any{"login": req.Login})
	if token == nil {
		return nil, nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error generating authentication token",
		}
	}

	return &GophermartUserRegisterResponse{
			Token: *token,
		}, &APIStatus{
			StatusCode: http.StatusOK,
			Message:    "User successfully registered and authenticated",
		}, nil
}

type GophermartUserRegisterRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type GophermartUserRegisterResponse struct {
	Token string `json:"token"`
}
