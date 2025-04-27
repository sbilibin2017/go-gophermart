package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type GophermartUserBalanceWithdrawalsListRepository interface {
	ListWithdrawalsByLoginOrdered(ctx context.Context, login string) ([]map[string]any, error)
}

type GophermartUserBalanceWithdrawalsListService struct {
	repo GophermartUserBalanceWithdrawalsListRepository
}

func NewGophermartUserBalanceWithdrawalsListService(
	repo GophermartUserBalanceWithdrawalsListRepository,
) *GophermartUserBalanceWithdrawalsListService {
	return &GophermartUserBalanceWithdrawalsListService{
		repo: repo,
	}
}

func (svc *GophermartUserBalanceWithdrawalsListService) List(
	ctx context.Context, login string,
) ([]*types.GophermartUserBalanceWithdrawalsResponse, *types.APIStatus, *types.APIStatus) {
	withdrawals, err := svc.repo.ListWithdrawalsByLoginOrdered(ctx, login)
	if err != nil {
		logger.Logger.Errorf("Error retrieving withdrawals for user %v: %v", login, err)
		return nil, nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error retrieving withdrawal data",
		}
	}
	if len(withdrawals) == 0 {
		return nil, &types.APIStatus{
			StatusCode: http.StatusNoContent,
			Message:    "No withdrawals found",
		}, nil
	}
	var response []*types.GophermartUserBalanceWithdrawalsResponse
	err = mapListToStruct(&response, withdrawals)
	if err != nil {
		logger.Logger.Errorf("Error mapping withdrawals: %v", err)
		return nil, nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error mapping withdrawal data",
		}
	}
	return response, &types.APIStatus{
		StatusCode: http.StatusOK,
		Message:    "Successfully retrieved withdrawals",
	}, nil
}
