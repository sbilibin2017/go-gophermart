package types

type OrderAccrual struct {
	Number  string `json:"number" db:"number"`
	Status  string `json:"status" db:"status"`
	Accrual *int64 `json:"accrual,omitempty" db:"accrual"`
}
