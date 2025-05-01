package types

type OrderRequest struct {
	Order string `json:"order" validate:"required,luhn"`
	Goods []Good `json:"goods" validate:"required,dive,required"`
}

type Good struct {
	Description string `json:"description" validate:"required"`
	Price       int64  `json:"price" validate:"required,gt=0"`
}

type OrderGetRequest struct {
	Number string `json:"number" validate:"required,luhn"`
}

type OrderGetResponse struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual int64  `json:"accrual"`
}

type OrderDB struct {
	Number  string `db:"number"`
	Status  string `db:"status"`
	Accrual int64  `db:"accrual"`
}

const (
	ORDER_ACCRUAL_STATUS_REGISTERED = "REGISTERED"
	ORDER_ACCRUAL_STATUS_INVALID    = "INVALID"
	ORDER_ACCRUAL_STATUS_PROCESSING = "PROCESSING"
	ORDER_ACCRUAL_STATUS_PROCESSED  = "PROCESSED"
)
