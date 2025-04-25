package models

import "time"

type RewardMechanicDB struct {
	Match      string    `db:"match"`
	Reward     int64     `db:"reward"`
	RewardType string    `db:"reward_type"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
