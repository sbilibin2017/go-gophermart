package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

type GophermartUserBalanceWithdrawService struct {
	v   StructValidator
	gbr FilterRepository
	sbr SaveRepository
}

func NewGophermartUserBalanceWithdrawService(
	v StructValidator,
	gbr FilterRepository,
	sbr SaveRepository,
) *GophermartUserBalanceWithdrawService {
	return &GophermartUserBalanceWithdrawService{
		v:   v,
		gbr: gbr,
		sbr: sbr,
	}
}

func (svc *GophermartUserBalanceWithdrawService) Withdraw(
	ctx context.Context, req *GophermartUserBalanceWithdrawRequest, login string,
) (*APIStatus, *APIStatus) {
	if err := svc.v.Struct(req); err != nil {
		return nil, &APIStatus{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Invalid order number format or invalid sum",
		}
	}
	balanceFilter := map[string]any{"login": login}
	balance, err := svc.gbr.Filter(ctx, balanceFilter, []string{"balance"})
	if err != nil {
		logger.Logger.Errorf("Error getting balance for user %v: %v", login, err)
		return nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error retrieving balance",
		}
	}
	currentBalance, ok := balance["balance"].(float64)
	if !ok {
		return nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error retrieving balance data",
		}
	}
	if currentBalance < float64(req.Sum) {
		return nil, &APIStatus{
			StatusCode: http.StatusPaymentRequired,
			Message:    "Insufficient balance",
		}
	}
	withdrawData := map[string]any{
		"login": login,
		"order": req.Order,
		"sum":   req.Sum,
	}
	err = svc.sbr.Save(ctx, withdrawData)
	if err != nil {
		logger.Logger.Errorf("Error saving withdrawal request for user %v: %v", login, err)
		return nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error saving withdrawal request",
		}
	}
	return &APIStatus{
		StatusCode: http.StatusOK,
		Message:    "Withdrawal request successfully registered",
	}, nil
}

type GophermartUserBalanceWithdrawRequest struct {
	Order string `json:"order" validate:"required,luhn"`
	Sum   int64  `json:"sum" validate:"required,gt=0"`
}
