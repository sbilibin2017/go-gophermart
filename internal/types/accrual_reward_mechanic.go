package types

import "time"

type AccrualRewardMechanicRegisterRequest struct {
	Match      string `json:"match"`
	Reward     int64  `json:"reward"`
	RewardType string `json:"reward_type"`
}

type AccrualRewardMechanicDB struct {
	Match      string    `db:"match"`
	Reward     int64     `db:"reward"`
	RewardType string    `db:"reward_type"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
