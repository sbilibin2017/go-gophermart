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
