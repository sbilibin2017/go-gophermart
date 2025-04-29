package services

import (
	"context"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserRegisterUserExistsByLoginRepository interface {
	ExistsByLogin(ctx context.Context, login string) (bool, error)
}

type UserRegisterUserSaveRepository interface {
	Save(ctx context.Context, login string, password string) error
}

type UserRegisterService struct {
	ueRepo UserRegisterUserExistsByLoginRepository
	usRepo UserRegisterUserSaveRepository
}

func NewUserRegisterService(ueRepo UserRegisterUserExistsByLoginRepository, usRepo UserRegisterUserSaveRepository) *UserRegisterService {
	return &UserRegisterService{
		ueRepo: ueRepo,
		usRepo: usRepo,
	}
}

func (s *UserRegisterService) Register(
	ctx context.Context, user *domain.User,
) error {
	exists, err := s.ueRepo.ExistsByLogin(ctx, user.Login)
	if err != nil {
		return err
	}
	if exists {
		return domain.ErrLoginAlreadyTaken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.usRepo.Save(ctx, user.Login, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}
