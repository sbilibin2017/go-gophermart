package models

import "time"

type GoodDB struct {
	Description string    `db:"description"`
	Price       float64   `db:"price"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time ` db:"updated_at"`
}
