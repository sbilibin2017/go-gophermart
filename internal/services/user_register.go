package services

import (
	"context"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserRegisterUserExistsByLoginRepository interface {
	ExistsByLogin(ctx context.Context, login string) (bool, error)
}

type UserRegisterUserSaveRepository interface {
	Save(ctx context.Context, user *domain.User) error
}

type UserRegisterService struct {
	ueRepo       UserRegisterUserExistsByLoginRepository
	usRepo       UserRegisterUserSaveRepository
	jwtSecretKey string
	jwtExp       time.Duration
}

func NewUserRegisterService(
	jwtSecretKey string,
	jwtExp time.Duration,
	ueRepo UserRegisterUserExistsByLoginRepository,
	usRepo UserRegisterUserSaveRepository,
) *UserRegisterService {
	return &UserRegisterService{
		ueRepo:       ueRepo,
		usRepo:       usRepo,
		jwtSecretKey: jwtSecretKey,
		jwtExp:       jwtExp,
	}
}

func (svc *UserRegisterService) Register(
	ctx context.Context, user *domain.User,
) (*string, error) {
	exists, err := svc.ueRepo.ExistsByLogin(ctx, user.Login)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrLoginAlreadyTaken
	}

	password, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password), bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}
	user.Password = string(password)

	err = svc.usRepo.Save(ctx, user)
	if err != nil {
		return nil, err
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
