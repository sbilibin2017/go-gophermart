package types

import "errors"

type GoodReward struct {
	Match      string `json:"match" db:"match" validate:"required"`
	Reward     int64  `json:"reward" db:"reward" validate:"gt=0"`
	RewardType string `json:"reward_type" db:"reward_type" validate:"oneof=% pt"`
}

var (
	ErrGoodRewardAlreadyExists = errors.New("good reward already exists")
)

const (
	REWARD_TYPE_PERCENT = "%"
	REWARD_TYPE_POINT   = "pt"
)
