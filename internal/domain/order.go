package domain

import (
	"errors"
	"time"
)

type Order struct {
	Number     string     `json:"number,omitempty" db:"number"`
	Login      string     `json:"login,omitempty" db:"login"`
	Status     string     `json:"status,omitempty" db:"status"`
	Accrual    *int64     `json:"accrual,omitempty" db:"accrual"`
	UploadedAt *time.Time `json:"uploaded_at,omitempty" db:"uploaded_at"`
}

var (
	ErrUserOrderExists = errors.New("user has already placed this order")
	ErrOrderExists     = errors.New("order with this number already exists")
)

const (
	GOPHERMART_ORDER_STATUS_NEW        = "NEW"
	GOPHERMART_ORDER_STATUS_PROCESSING = "PROCESSING"
	GOPHERMART_ORDER_STATUS_INVALID    = "INVALID"
	GOPHERMART_ORDER_STATUS_PROCESSED  = "PROCESSED"
)
