package types

type RegisterOrderRequest struct {
	Order string              `json:"order" validate:"required,luhn"`
	Goods []RegisterOrderGood `json:"goods" validate:"required"`
}

type RegisterOrderGood struct {
	Description string `json:"description" validate:"required"`
	Price       int64  `json:"price" validate:"required,gt=0"`
}

type OrderStatus string

const (
	OrderStatusRegistered OrderStatus = "REGISTERED"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type OrderResponse struct {
	Order   string      `json:"order"`
	Status  OrderStatus `json:"status"`
	Accrual *int64      `json:"accrual,omitempty"`
}

func NewOrderResponse(result map[string]any) *OrderResponse {
	response := &OrderResponse{
		Order:  result["order_id"].(string),
		Status: OrderStatus(result["status"].(string)),
	}
	if accrualRaw, ok := result["accrual"]; ok && accrualRaw != nil {
		if accrual, ok := accrualRaw.(int64); ok {
			response.Accrual = &accrual
		}
	}
	return response
}
