package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type GophermartUserBalanceGetBalanceRepository interface {
	GetBalance(ctx context.Context, login string) (float64, error)
}

type GophermartUserBalanceGetBalanceWithdrawalSaveRepository interface {
	Save(ctx context.Context, login, order string, sum int64) error
}

type GophermartUserBalanceWithdrawValidator interface {
	Struct(v any) error
}

type GophermartUserBalanceWithdrawService struct {
	v   GophermartUserBalanceWithdrawValidator
	gbr GophermartUserBalanceGetBalanceRepository
	sbr GophermartUserBalanceGetBalanceWithdrawalSaveRepository
}

func NewGophermartUserBalanceWithdrawService(
	v GophermartUserBalanceWithdrawValidator,
	gbr GophermartUserBalanceGetBalanceRepository,
	sbr GophermartUserBalanceGetBalanceWithdrawalSaveRepository,
) *GophermartUserBalanceWithdrawService {
	return &GophermartUserBalanceWithdrawService{
		v:   v,
		gbr: gbr,
		sbr: sbr,
	}
}

func (svc *GophermartUserBalanceWithdrawService) Withdraw(
	ctx context.Context, req *types.GophermartUserBalanceWithdrawRequest, login string,
) (*types.APIStatus, *types.APIStatus) {
	if err := svc.v.Struct(req); err != nil {
		return nil, &types.APIStatus{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Invalid order number format or invalid sum",
		}
	}
	currentBalance, err := svc.gbr.GetBalance(ctx, login)
	if err != nil {
		logger.Logger.Errorf("Error getting balance for user %v: %v", login, err)
		return nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error retrieving balance",
		}
	}
	if currentBalance < float64(req.Sum) {
		return nil, &types.APIStatus{
			StatusCode: http.StatusPaymentRequired,
			Message:    "Insufficient balance",
		}
	}
	err = svc.sbr.Save(ctx, login, req.Order, req.Sum)
	if err != nil {
		logger.Logger.Errorf("Error saving withdrawal request for user %v: %v", login, err)
		return nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error saving withdrawal request",
		}
	}
	return &types.APIStatus{
		StatusCode: http.StatusOK,
		Message:    "Withdrawal request successfully registered",
	}, nil
}
