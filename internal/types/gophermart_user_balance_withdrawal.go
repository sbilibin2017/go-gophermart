package types

import "time"

type GophermartUserBalanceWithdrawRequest struct {
	Order string `json:"order" validate:"required,luhn"`
	Sum   int64  `json:"sum" validate:"required,gt=0"`
}

type GophermartUserBalanceWithdrawalsResponse struct {
	Order       string    `json:"order"`
	Sum         int64     `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type GophermartUserBalanceWithdrawalDB struct {
	Login       string    `db:"login"`
	Number      string    `db:"number"`
	Sum         int64     `db:"sum"`
	ProcessedAt time.Time `db:"processed_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
