package services

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/jwt"
)

type UserLoginUserGetByLoginRepository interface {
	GetByLogin(ctx context.Context, login string) (map[string]any, error)
}

type UserLoginService struct {
	gblRepo      UserLoginUserGetByLoginRepository
	jwtSecretKey string
	jwtExp       time.Duration
}

func NewUserLoginService(
	jwtSecretKey string,
	jwtExp time.Duration,
	gblRepo UserLoginUserGetByLoginRepository,
) *UserLoginService {
	return &UserLoginService{
		gblRepo:      gblRepo,
		jwtSecretKey: jwtSecretKey,
		jwtExp:       jwtExp,
	}
}

func (svc *UserLoginService) Login(
	ctx context.Context, user *domain.User,
) (*string, error) {
	u, err := svc.gblRepo.GetByLogin(ctx, user.Login)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(u["password"].(string)),
		[]byte(user.Password),
	)
	if err != nil {
		return nil, domain.ErrInvalidUserCredentials
	}
	token, err := jwt.GenerateTokenString(
		user.Login,
		svc.jwtSecretKey,
		svc.jwtExp,
	)
	if err != nil {
		return nil, err
	}
	return &token, nil
}
