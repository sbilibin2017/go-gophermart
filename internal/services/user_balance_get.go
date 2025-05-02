package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserBalanceGetUserBalanceFilterOneRepository interface {
	FilterOne(ctx context.Context, login string) (*types.UserBalanceDB, error)
}

type UserBalanceGetService struct {
	ubg UserBalanceGetUserBalanceFilterOneRepository
}

func NewUserBalanceGetService(
	ubg UserBalanceGetUserBalanceFilterOneRepository,
) *UserBalanceGetService {
	return &UserBalanceGetService{
		ubg: ubg,
	}
}

func (svc *UserBalanceGetService) Get(
	ctx context.Context, login string,
) (*types.UserBalanceResponse, *types.APIStatus) {
	userBalance, err := svc.ubg.FilterOne(ctx, login)
	if err != nil {
		return nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch user balance",
		}
	}
	return &types.UserBalanceResponse{
		Current:   userBalance.Current,
		Withdrawn: userBalance.Withdrawn,
	}, nil
}
