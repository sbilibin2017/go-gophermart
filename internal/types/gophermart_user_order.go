package types

import "time"

type GophermartUserOrder struct {
	Login      string    `db:"login"`
	Number     string    `db:"number"`
	Status     string    `db:"status"`
	Accrual    int64     `db:"accrual"`
	UploadedAt time.Time `db:"uploaded_at"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
