package types

import "errors"

type OrderRequest struct {
	Order string `json:"order" validate:"required,luhn"`
	Goods []Good `json:"goods" validate:"required,dive"`
}

type Good struct {
	Description string `json:"description" validate:"required"`
	Price       int64  `json:"price" validate:"required,gt=0"`
}

type Order struct {
	Number  string `json:"number" db:"number"`
	Accrual *int64 `json:"accrual,omitempty" db:"accrual,omitempty"`
	Status  string `json:"status" db:"status"`
}

const (
	ORDER_STATUS_REGISTERED = "REGISTERED"
	ORDER_STATUS_INVALID    = "INVALID"
	ORDER_STATUS_PROCESSING = "PROCESSING"
	ORDER_STATUS_PROCESSED  = "PROCESSED"
)

var (
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrOrderNotFound      = errors.New("order not found")
)
