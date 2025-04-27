package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type GophermartUserCurrentBalanceRepository interface {
	GetBalanceByLogin(ctx context.Context, login string, fields []string) (map[string]any, error)
}

type GophermartUserCurrentBalanceService struct {
	ro GophermartUserCurrentBalanceRepository
}

func NewGophermartUserCurrentBalanceService(
	ro GophermartUserCurrentBalanceRepository,
) *GophermartUserCurrentBalanceService {
	return &GophermartUserCurrentBalanceService{
		ro: ro,
	}
}

func (svc *GophermartUserCurrentBalanceService) Get(
	ctx context.Context, login string,
) (*types.GophermartUserCurrentBalanceResponse, *types.APIStatus, *types.APIStatus) {
	balance, err := svc.ro.GetBalanceByLogin(ctx, login, []string{"current", "withdrawn"})
	if err != nil {
		logger.Logger.Errorf("Error getting balance for user %v: %v", login, err)
		return nil, nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error retrieving balance",
		}
	}
	var response *types.GophermartUserCurrentBalanceResponse
	err = mapToStruct(response, balance)
	if err != nil {
		return nil, nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error mapping balance data",
		}
	}
	return response, &types.APIStatus{
		StatusCode: http.StatusOK,
		Message:    "Successfully retrieved balance",
	}, nil
}
