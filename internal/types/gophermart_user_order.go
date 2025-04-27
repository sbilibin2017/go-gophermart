package types

import "time"

type GophermartUserOrderUploadRequest struct {
	Number string `json:"number" validate:"required,luhn"`
}

type GophermartUserOrdersResponse struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    *int64    `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type GophermartUserOrderDB struct {
	Login      string    `db:"login"`
	Number     string    `db:"number"`
	Status     string    `db:"status"`
	Accrual    int64     `db:"accrual"`
	UploadedAt time.Time `db:"uploaded_at"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
