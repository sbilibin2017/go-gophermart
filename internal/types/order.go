package types

type OrderResponse struct {
	Order   string      `json:"order" `
	Status  OrderStatus `json:"status"`
	Accrual *int64      `json:"accrual"`
}

type OrderStatus string

const (
	OrderStatusRegistered OrderStatus = "REGISTERED"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)
