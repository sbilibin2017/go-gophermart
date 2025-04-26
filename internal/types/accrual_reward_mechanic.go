package types

import "time"

type AccrualRewardMechanic struct {
	Match      string    `db:"match"`
	Reward     int64     `db:"reward"`
	RewardType string    `db:"reward_type"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
