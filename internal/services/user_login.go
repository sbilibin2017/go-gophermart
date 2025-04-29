package services

import (
	"context"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
)

type UserLoginUserGetByLoginRepository interface {
	GetByLogin(ctx context.Context, login string) (*domain.User, error)
}

type UserLoginService struct {
	gblRepo          UserLoginUserGetByLoginRepository
	jwtSecretKey     string
	jwtExp           time.Duration
	jwtGenerator     func(login string, secretKey string, exp time.Duration) (string, error)
	passwordComparer func(hashedPassword []byte, password []byte) error
}

func NewUserLoginService(
	jwtSecretKey string,
	jwtExp time.Duration,
	jwtGenerator func(login string, secretKey string, exp time.Duration) (string, error),
	passwordComparer func(hashedPassword []byte, password []byte) error,
	gblRepo UserLoginUserGetByLoginRepository,
) *UserLoginService {
	return &UserLoginService{
		gblRepo:          gblRepo,
		jwtSecretKey:     jwtSecretKey,
		jwtExp:           jwtExp,
		jwtGenerator:     jwtGenerator,
		passwordComparer: passwordComparer,
	}
}

func (svc *UserLoginService) Login(
	ctx context.Context, user *domain.User,
) (*string, error) {
	u, err := svc.gblRepo.GetByLogin(ctx, user.Login)
	if err != nil {
		return nil, err
	}
	err = svc.passwordComparer(
		[]byte(u.Password),
		[]byte(user.Password),
	)
	if err != nil {
		return nil, domain.ErrInvalidUserCredentials
	}
	token, err := svc.jwtGenerator(
		user.Login,
		svc.jwtSecretKey,
		svc.jwtExp,
	)
	if err != nil {
		return nil, err
	}
	return &token, nil
}
