package types

import "time"

type OrderGetRequest struct {
	Number string `json:"number" validate:"required,luhn"`
}

type OrderGetResponse struct {
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
	OrderStatusRegistered OrderStatus = "REGISTERED" // Заказ зарегистрирован, но начисление не рассчитано
	OrderStatusInvalid    OrderStatus = "INVALID"    // Заказ не принят к расчёту, вознаграждение не будет начислено
	OrderStatusProcessing OrderStatus = "PROCESSING" // Расчёт начисления в процессе
	OrderStatusProcessed  OrderStatus = "PROCESSED"  // Расчёт начисления окончен
)
