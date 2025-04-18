package domain

import "errors"

type Reward struct {
	Match      string
	Reward     uint64
	RewardType RewardType
}

func NewReward(
	match string,
	reward uint64,
	rewardType string,
) *Reward {
	return &Reward{
		Match:      match,
		Reward:     reward,
		RewardType: RewardType(rewardType),
	}
}

type RewardType string

const (
	RewardTypePercent RewardType = "%"
	RewardTypePoints  RewardType = "pt"
)

var ErrRewardSearchKeyAlreadyRegistered = errors.New("reward search key already registered")
