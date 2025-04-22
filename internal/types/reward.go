package types

import "time"

type RewardRegisterRequest struct {
	Match      string     `json:"match" validate:"required"`
	Reward     int64      `json:"reward" validate:"required,gt=0"`
	RewardType RewardType `json:"reward_type" validate:"required,reward_type"`
}

type RewardDB struct {
	Match      string     `db:"match"`
	Reward     int64      `db:"reward"`
	RewardType RewardType `db:"reward_type"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
}

type RewardType string

const (
	RewardTypePercent RewardType = "%"
	RewardTypePoint   RewardType = "pt"
)
