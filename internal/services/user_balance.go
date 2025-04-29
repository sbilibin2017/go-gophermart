package services

import (
	"context"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
)

type UserBalanceRepository interface {
	GetByLogin(ctx context.Context, login string) (*domain.UserBalance, error)
}

type UserBalanceService struct {
	balanceRepo UserBalanceRepository
}

func NewUserBalanceService(
	balanceRepo UserBalanceRepository,
) *UserBalanceService {
	return &UserBalanceService{
		balanceRepo: balanceRepo,
	}
}

func (svc *UserBalanceService) GetBalance(
	ctx context.Context, login string,
) (*domain.UserBalance, error) {
	balance, err := svc.balanceRepo.GetByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	return balance, nil
}
