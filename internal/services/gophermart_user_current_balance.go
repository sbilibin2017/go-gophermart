package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

type GophermartUserCurrentBalanceService struct {
	ro FilterRepository
}

func NewGophermartUserCurrentBalanceService(
	ro FilterRepository,
) *GophermartUserCurrentBalanceService {
	return &GophermartUserCurrentBalanceService{
		ro: ro,
	}
}

func (svc *GophermartUserCurrentBalanceService) Get(
	ctx context.Context, login string,
) (*GophermartUserCurrentBalanceResponse, *APIStatus, *APIStatus) {
	filter := map[string]any{"login": login}
	balance, err := svc.ro.Filter(ctx, filter, []string{"current", "withdrawn"})
	if err != nil {
		logger.Logger.Errorf("Error getting balance for user %v: %v", login, err)
		return nil, nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error retrieving balance",
		}
	}
	var response *GophermartUserCurrentBalanceResponse
	err = mapToStruct(response, balance)
	if err != nil {
		return nil, nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error mapping balance data",
		}
	}
	return response, &APIStatus{
		StatusCode: http.StatusOK,
		Message:    "Successfully retrieved balance",
	}, nil
}

type GophermartUserCurrentBalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn int64   `json:"withdrawn"`
}
