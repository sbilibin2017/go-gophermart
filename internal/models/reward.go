package models

import "time"

type RewardDB struct {
	Match      string    `db:"match"`
	Reward     uint64    `db:"reward"`
	RewardType string    `db:"reward_type"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type RewardFilter struct {
	Match string `db:"match"`
}
