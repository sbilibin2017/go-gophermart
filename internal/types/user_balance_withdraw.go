package types

import "time"

type UserBalanceWithdrawRequest struct {
	Login string `json:"login"`
	Order string `json:"order" validate:"required,luhn"`
	Sum   int64  `json:"sum" validate:"required,ge=0"`
}

type UserWithdrawResponse struct {
	Login       string    `json:"login"`
	Number      string    `json:"number"`
	Sum         int64     `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type UserBalanceWithdrawDB struct {
	Login       string    `db:"login"`
	Number      string    `db:"number"`
	Sum         int64     `db:"sum"`
	ProcessedAt time.Time `db:"processed_at"`
}
