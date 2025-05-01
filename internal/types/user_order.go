package types

import "time"

type UserOrderUploadRequest struct {
	Login  string `json:"login"`
	Number string `json:"number" validate:"required,luhn"`
}

type UserOrdersRequest string

type UserOrdersResponse struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    *int64    `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type UserOrderDB struct {
	Login      string    `db:"login"`
	Number     string    `db:"number"`
	Status     string    `db:"status"`
	Accrual    *int64    `db:"accrual"`
	UploadedAt time.Time `db:"uploaded_at"`
}

const (
	ORDER_GOPHERMART_STATUS_REGISTERED = "NEW"
	ORDER_GOPHERMART_STATUS_INVALID    = "PROCESSING"
	ORDER_GOPHERMART_STATUS_PROCESSING = "INVALID"
	ORDER_GOPHERMART_STATUS_PROCESSED  = "PROCESSED"
)
