package models

import "time"

type OrderDB struct {
	Number    string    `db:"number"`
	Login     int       `db:"login"`
	Status    string    `db:"status"`
	Accrual   *int64    `db:"accrual"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
