package types

import "time"

type UserOrderUploadNumberRequest struct {
	Login  string `json:"login"`
	Number string `json:"number" validate:"required,luhn"`
}

type UserOrderUploadedListResponse struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    *int64    `json:"accural,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type UserOrderDB struct {
	Number     string    `db:"number"`
	Login      string    `db:"login"`
	Status     string    `db:"status"`
	Accrual    *int64    `db:"accrual"`
	UploadedAt time.Time `db:"uploaded_at"`
}

const (
	GophermartUserOrderStatusNew        = "NEW"
	GophermartUserOrderStatusProcessing = "PROCESSING"
	GophermartUserOrderStatusInvalid    = "INVALID"
	GophermartUserOrderStatusProcessed  = "PROCESSED"
)
