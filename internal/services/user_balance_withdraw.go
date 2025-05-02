package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/services/validation"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserBalanceWithdrawUserBalanceFilterOneRepository interface {
	FilterOne(ctx context.Context, login string) (*types.UserBalanceDB, error)
}

type UserBalanceWithdrawSaveRepository interface {
	Save(ctx context.Context, login string, number string, sum int64) error
}

type UserBalanceWithdrawValidator interface {
	Struct(v any) error
}

type UserBalanceWithdrawService struct {
	v    UserBalanceWithdrawValidator
	ubf  UserBalanceWithdrawUserBalanceFilterOneRepository
	ubws UserBalanceWithdrawSaveRepository
}

func NewUserBalanceWithdrawService(
	v UserBalanceWithdrawValidator,
	ubf UserBalanceWithdrawUserBalanceFilterOneRepository,
	ubws UserBalanceWithdrawSaveRepository,
) *UserBalanceWithdrawService {
	return &UserBalanceWithdrawService{
		v:    v,
		ubf:  ubf,
		ubws: ubws,
	}
}

func (svc *UserBalanceWithdrawService) Withdraw(
	ctx context.Context, req *types.UserBalanceWithdrawRequest,
) *types.APIStatus {
	if err := svc.v.Struct(req); err != nil {
		valErr := formatValidationError(err)
		if valErr != nil {
			if validation.IsLuhnValidationError(valErr) {
				return &types.APIStatus{
					StatusCode: http.StatusUnprocessableEntity,
					Message:    valErr.Message,
				}
			}
			return &types.APIStatus{
				StatusCode: http.StatusBadRequest,
				Message:    valErr.Message,
			}
		}
		return &types.APIStatus{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request structure",
		}
	}

	balance, err := svc.ubf.FilterOne(ctx, req.Login)
	if err != nil {
		return &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get user balance",
		}
	}
	if balance.Current < float64(req.Sum) {
		return &types.APIStatus{
			StatusCode: http.StatusPaymentRequired,
			Message:    "Insufficient funds",
		}
	}

	if err := svc.ubws.Save(ctx, req.Login, req.Order, req.Sum); err != nil {
		return &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to register withdrawal",
		}
	}

	return &types.APIStatus{
		StatusCode: http.StatusOK,
		Message:    "Withdrawal successful",
	}
}
