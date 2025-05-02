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
	GOPHERMART_USER_ORDER_STATUS_NEW        = "NEW"
	GOPHERMART_USER_ORDER_STATUS_PROCESSING = "PROCESSING"
	GOPHERMART_USER_ORDER_STATUS_INVALID    = "INVALID"
	GOPHERMART_USER_ORDER_STATUS_PROCESSED  = "PROCESSED"
)
