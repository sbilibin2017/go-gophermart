package types

import "errors"

type UserOrder struct {
	Number     string `json:"number" db:"number"`
	Login      string `json:"login" db:"login"`
	Status     string `json:"status" db:"status"`
	Accrual    *int64 `json:"accrual,omitempty" db:"accrual"`
	UploadedAt string `json:"uploaded_at" db:"uploaded_at"`
}

var (
	ErrOrderAlreadyUploadedByUser      = errors.New("order number already uploaded by this user")
	ErrOrderAlreadyUploadedByOtherUser = errors.New("order number already uploaded by another user")
)
