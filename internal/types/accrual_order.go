package types

import "time"

type AccrualOrder struct {
	Number    string    `db:"number"`
	Status    string    `db:"status"`
	Accrual   int64     `db:"accrual"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
