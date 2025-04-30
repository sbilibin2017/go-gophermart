package types

type Order struct {
	Number     string `json:"number" db:"number"`
	Login      string `json:"login" db:"login"`
	Status     string `json:"status" db:"status"`
	Accrual    *int64 `json:"accrual,omitempty" db:"accrual"`
	UploadedAt string `json:"uploaded_at" db:"uploaded_at"`
}
