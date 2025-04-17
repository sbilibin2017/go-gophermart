package domain

import "errors"

type Reward struct {
	Match      string     `json:"match" validate:"required"`
	Reward     uint       `json:"reward" validate:"required,gt=0"`
	RewardType RewardType `json:"reward_type" validate:"required,oneof=% pt"`
}

type RewardType string

const (
	RewardTypePercent RewardType = "%"
	RewardTypePoints  RewardType = "pt"
)

var (
	ErrRewardKeyAlreadyRegistered = errors.New("reward key already registered")
	ErrFailedToRegisterReward     = errors.New("failed to register reward")
)
