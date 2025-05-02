package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserBalanceWithdrawListRepository interface {
	ListOrdered(ctx context.Context, login string) (*[]types.UserBalanceWithdrawDB, error)
}

type UserBalanceWithdrawListService struct {
	ubwRepo UserBalanceWithdrawListRepository
}

func NewUserBalanceWithdrawListService(
	ubwRepo UserBalanceWithdrawListRepository,
) *UserBalanceWithdrawListService {
	return &UserBalanceWithdrawListService{
		ubwRepo: ubwRepo,
	}
}

func (svc *UserBalanceWithdrawListService) List(
	ctx context.Context, login string,
) ([]*types.UserBalanceWithdrawResponse, *types.APIStatus) {
	withdrawals, err := svc.ubwRepo.ListOrdered(ctx, login)
	if err != nil {
		return nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch user withdrawals",
		}
	}

	if len(*withdrawals) == 0 {
		return nil, &types.APIStatus{
			StatusCode: http.StatusNoContent,
			Message:    "No withdrawals found",
		}
	}

	var response []*types.UserBalanceWithdrawResponse
	for _, withdrawal := range *withdrawals {
		response = append(response, &types.UserBalanceWithdrawResponse{
			Order:       withdrawal.Number,
			Sum:         withdrawal.Sum,
			ProcessedAt: withdrawal.ProcessedAt,
		})
	}

	return response, &types.APIStatus{
		StatusCode: http.StatusOK,
		Message:    "Withdrawals fetched successfully",
	}
}
