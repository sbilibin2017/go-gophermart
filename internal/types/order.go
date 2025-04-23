package types

import "time"

type OrderGetByIDRequest struct {
	Number string `json:"number" validate:"required,luhn"`
}

type OrderGetByIDResponse struct {
	Order   string      `json:"order"`
	Status  OrderStatus `json:"status"`
	Accrual *int64      `json:"accrual,omitempty"`
}

type OrderAcceptRequest struct {
	Order string `json:"order" validate:"required,luhn"`
	Goods []Good `json:"goods" validate:"required,min=1"`
}

type Good struct {
	Description string `json:"description" validate:"required"`
	Price       int64  `json:"price" validate:"required,gt=0"`
}

type OrderDB struct {
	OrderID   string      `db:"order_id"`
	Status    OrderStatus `db:"status"`
	Accrual   *int64      `db:"accrual"`
	CreatedAt time.Time   `db:"created_at"`
	UpdatedAt time.Time   `db:"updated_at"`
}

type OrderStatus string

const (
	OrderStatusRegistered OrderStatus = "REGISTERED"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)
