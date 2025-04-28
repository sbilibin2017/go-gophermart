package services

import (
	"context"
	"net/http"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

type GophermartUserBalanceWithdrawalsListService struct {
	repo ListRepository
}

func NewGophermartUserBalanceWithdrawalsListService(
	repo ListRepository,
) *GophermartUserBalanceWithdrawalsListService {
	return &GophermartUserBalanceWithdrawalsListService{
		repo: repo,
	}
}

func (svc *GophermartUserBalanceWithdrawalsListService) List(
	ctx context.Context, login string,
) ([]*GophermartUserBalanceWithdrawalsResponse, *APIStatus, *APIStatus) {
	filter := map[string]any{"login": login}
	withdrawals, err := svc.repo.List(ctx, filter, []string{"order", "sum", "processed_at"}, "processed_at", true)
	if err != nil {
		logger.Logger.Errorf("Error retrieving withdrawals for user %v: %v", login, err)
		return nil, nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error retrieving withdrawal data",
		}
	}
	if len(withdrawals) == 0 {
		return nil, &APIStatus{
			StatusCode: http.StatusNoContent,
			Message:    "No withdrawals found",
		}, nil
	}
	var response []*GophermartUserBalanceWithdrawalsResponse
	err = mapListToStruct(&response, withdrawals)
	if err != nil {
		logger.Logger.Errorf("Error mapping withdrawals: %v", err)
		return nil, nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error mapping withdrawal data",
		}
	}
	return response, &APIStatus{
		StatusCode: http.StatusOK,
		Message:    "Successfully retrieved withdrawals",
	}, nil
}

type GophermartUserBalanceWithdrawalsResponse struct {
	Order       string    `json:"order"`
	Sum         int64     `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
